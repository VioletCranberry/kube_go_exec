package main

import (
	"kube_go_exec/utils"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var logger = logrus.New()
var config = viper.New()

func main() {

	logger.Formatter = &logrus.JSONFormatter{}

	config.AddConfigPath(".")
	config.SetConfigName("app")
	config.SetConfigType("env")
	config.AutomaticEnv()
	err := config.ReadInConfig()
	if err != nil {
		logger.Fatal(err)
	}

	k8s_client, err := utils.KubeClientFromConfig(
		config.GetString("KUBERNETES_BEARER_CONFIG"),
		config.GetString("KUBERNETES_MASTER_URL"),
	)
	if err != nil {
		logger.Fatal(err)
	}

	pod, err := k8s_client.GetPodByFilter(
		config.GetString("KUBERNETES_POD_LABEL"),
		config.GetString("KUBERNETES_NAMESPACE"),
	)
	if err != nil {
		logger.Fatal(err)
	}

	logger.WithFields(logrus.Fields{
		"pod":        pod.Name,
		"namespace":  pod.Namespace,
		"status":     pod.Status.Phase,
		"labels":     pod.Labels,
		"containers": pod.Spec.Containers,
	}).Info("Pod Context")

	stdout, stderr, err := k8s_client.ExecInPod(pod,
		config.GetString("KUBERNETES_CONTAINER"),
		config.GetString("KUBERNETES_POD_EXEC"),
	)
	if err != nil {
		logger.Fatal(err)
	}

	slack_client := utils.InitSlackApi(
		config.GetString("SLACK_TOKEN"),
	)

	channelId, attachment, err := slack_client.SendToChannel(
		config.GetString("SLACK_CHANNEL_ID"),
		config.GetString("KUBERNETES_POD_EXEC"),
		pod.Name,
		stdout,
		stderr,
	)
	if err != nil {
		logger.Fatal(err)
	}

	logger.WithFields(logrus.Fields{
		"channel id": channelId,
		"attachment": attachment,
	}).Info("Slack Context")
}
