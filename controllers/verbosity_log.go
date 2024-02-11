package controllers

import (
	"context"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type VerbosityLog struct {
	ctx      context.Context
	l        logr.Logger
	level    int
	maxLevel int
}

func (x *VerbosityLog) FromContext(ctx context.Context) VerbosityLog {
	y := *x
	y.ctx = ctx
	y.l = log.FromContext(ctx)
	return y
}

func (x *VerbosityLog) SetMaxLevel(maxLevel int) VerbosityLog {
	x.maxLevel = maxLevel
	y := *x
	return y
}
func (x VerbosityLog) V(level int) VerbosityLog {
	y := x
	y.maxLevel = x.maxLevel
	y.l = x.l.V(level)
	y.level = level
	return y
}

func (x VerbosityLog) Error(err error, msg string, keysAndValues ...interface{}) {
	if x.l.Enabled() && x.level <= x.maxLevel {
		x.l.Error(err, msg, keysAndValues...)
	}
}

func (x VerbosityLog) Info(msg string, keysAndValues ...interface{}) {
	if x.l.Enabled() && x.level <= x.maxLevel {
		x.l.Info(msg, keysAndValues...)
	}
}
