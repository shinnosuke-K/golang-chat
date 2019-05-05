package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
)

// ErrNoAvatarはAvatarインスタンスがアバターのURLを返すことができる
// 場合に発生するエラーです
var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得できません。")

// Avatarはユーザーのプロフィール画像を表す型です。
type Avatar interface {
	// GetAvatarURLは指定されたクライアントのアバターのURLを返します
	// 問題が発生した場合にはエラーを返します。特に、URLを取得できなかった
	// 場合にはErrNoAvatarURLを返します
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (_ AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvararURL()
	if url != "" {
		return url, nil
	}
	return "", ErrNoAvatarURL
}

type GravatarAvatar struct{}

var UserGravatar GravatarAvatar

func (_ GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

type FileSystemAvatar struct{}

var UserFileSystemAvatar FileSystemAvatar

func (_ FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	if files, err := ioutil.ReadDir("avatars"); err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if match, _ := filepath.Match(u.UniqueID()+"*", file.Name()); match {
				return "/avatars/" + file.Name(), nil
			}
		}
	}
	return "", ErrNoAvatarURL
}
