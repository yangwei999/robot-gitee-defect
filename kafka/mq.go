package kafka

import (
	libkafka "github.com/opensourceways/kafka-lib/kafka"
	"github.com/opensourceways/kafka-lib/mq"
	"github.com/sirupsen/logrus"
)

var instance *serviceImpl

func Init(cfg *Config, log *logrus.Entry) error {
	err := libkafka.Init(
		mq.Addresses(cfg.mqConfig().Addresses...),
		mq.Log(log),
	)
	if err != nil {
		return err
	}

	if err := libkafka.Connect(); err != nil {
		return err
	}

	instance = &serviceImpl{}

	return nil
}

func Exit() {
	instance.unsubscribe()

	if err := libkafka.Disconnect(); err != nil {
		logrus.Errorf("exit kafka, err:%v", err)
	}
}

func Instance() *serviceImpl {
	return instance
}

type Handler func([]byte) error

type serviceImpl struct {
	subscribers []mq.Subscriber
}

func (impl *serviceImpl) Publish(topic string, msg []byte) error {
	return libkafka.Publish(topic, &mq.Message{
		Body: msg,
	})
}

func (impl *serviceImpl) unsubscribe() {
	s := impl.subscribers
	for i := range s {
		if err := s[i].Unsubscribe(); err != nil {
			logrus.Errorf(
				"failed to unsubscribe to topic:%s, err:%v",
				s[i].Topic, err,
			)
		}
	}
}

func (impl *serviceImpl) Subscribe(group string, handlers map[string]Handler) error {
	for topic, h := range handlers {
		s, err := impl.registerHandler(topic, group, h)
		if err != nil {
			return err
		}

		if s != nil {
			impl.subscribers = append(impl.subscribers, s)
		} else {
			logrus.Infof("does not subscribe topic:%s", topic)
		}
	}

	return nil
}

func (impl *serviceImpl) registerHandler(topic, group string, h Handler) (mq.Subscriber, error) {
	if h == nil {
		return nil, nil
	}

	return libkafka.Subscribe(topic, group+topic, func(e mq.Event) error {
		msg := e.Message()
		if msg == nil {
			return nil
		}

		return h(msg.Body)
	})
}
