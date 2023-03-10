package router

import (
	"social_network/internal/api/router/options"
	"social_network/internal/oauth/github"
)

func init(){
	router.APIRouter.HandleFunc("/login/github",github.GithubLogin)

	router.APIRouter.HandleFunc("/github/callback",github.GithubCallback)
}