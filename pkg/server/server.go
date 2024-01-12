package server

import (
	"k8s.io/test-infra/prow/github"
)

type GithubClient interface {
	GetPullRequest(org, repo string, num int) (*github.PullRequest, error)
	CreateComment(owner, repo string, number int, comment string) error
	AddLabel(owber, repo string, number int, label string) error
	RemoveLabel(owner, repo string, number int, label string) error
	WasLabelAddedByHuman(org, repo string, num int, label string) (bool, error)
	BotUserChecker() (func(candidate string) bool, error)
}
