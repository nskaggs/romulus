// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package cmd

import (
	"github.com/juju/cmd"
	"github.com/juju/errors"
	"github.com/juju/idmclient/ussologin"
	"github.com/juju/juju/jujuclient"
	"github.com/juju/persistent-cookiejar"
	"gopkg.in/juju/environschema.v1/form"
	"gopkg.in/macaroon-bakery.v1/httpbakery"
)

// TODO (mattyw) Http needs to be HTTP
// HttpCommand can instantiate http bakery clients using a common cookie jar.
type HttpCommand struct {
	cmd.CommandBase

	cookiejar *cookiejar.Jar
}

// NewClient returns a new HTTP bakery client for commands.
func (s *HttpCommand) NewClient(ctx *cmd.Context) (*httpbakery.Client, error) {
	if s.cookiejar == nil {
		cookieFile := cookiejar.DefaultCookieFile()
		jar, err := cookiejar.New(&cookiejar.Options{
			Filename: cookieFile,
		})
		if err != nil {
			return nil, errors.Trace(err)
		}
		s.cookiejar = jar
	}
	client := httpbakery.NewClient()
	client.Jar = s.cookiejar
	client.VisitWebPage = ussologin.VisitWebPage(newIOFiller(ctx), client.Client, jujuclient.NewTokenStore())
	return client, nil
}

// Close saves the persistent cookie jar used by the specified httpbakery.Client.
func (s *HttpCommand) Close() error {
	if s.cookiejar != nil {
		return s.cookiejar.Save()
	}
	return nil
}

func newIOFiller(ctx *cmd.Context) *form.IOFiller {
	return &form.IOFiller{
		In:  ctx.Stdin,
		Out: ctx.Stderr,
	}
}
