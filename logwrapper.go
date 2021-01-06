//    Copyright 2021 Stefen Sharkey
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type Event struct {
	level   logrus.Level
	message interface{}
}

type StandardLogger struct {
	*logrus.Logger
}

func NewLogger() *StandardLogger {
	baseLogger := logrus.New()
	standardLogger := &StandardLogger{baseLogger}
	standardLogger.Formatter = &logrus.JSONFormatter{}

	return standardLogger
}

var (
	joinedGuildMessage            = Event{logrus.InfoLevel, "Joined server: %s"}
	leftGuildMessage              = Event{logrus.InfoLevel, "Left server: %s"}
	startedMessage                = Event{logrus.InfoLevel, "TicketBot is now running. Press CTRL-C to exit."}
	stoppingMessage               = Event{logrus.InfoLevel, "TicketBot is now stopping."}
	configErrorMessage            = Event{logrus.ErrorLevel, "Config Error: %s"}
	discordConnectionErrorMessage = Event{logrus.ErrorLevel, "Discord Connection Error: %s"}
	discordSessionErrorMessage    = Event{logrus.ErrorLevel, "Discord Session Error: %s"}
	sqlErrorMessage               = Event{logrus.PanicLevel, "SQL Error: %s"}
	dsnDebugMessage               = Event{logrus.DebugLevel, "DSN: %s"}
	sqlOpeningMessage             = Event{logrus.TraceLevel, "Starting SQL Connection: %s"}
	sqlOpenedMessage              = Event{logrus.TraceLevel, "Started SQL Connection: %s"}
)

func (l *StandardLogger) JoinedGuild(guildName string) {
	l.HandleEvent(joinedGuildMessage, guildName)
}

func (l *StandardLogger) LeftGuild(guildName string) {
	l.HandleEvent(leftGuildMessage, guildName)
}

func (l *StandardLogger) Started() {
	l.HandleEvent(startedMessage)
}

func (l *StandardLogger) Stopping() {
	l.HandleEvent(stoppingMessage)
}

func (l *StandardLogger) ConfigError(errorMessage error) {
	l.HandleEvent(configErrorMessage, errorMessage)
}

func (l *StandardLogger) DiscordConnectionError(errorMessage error) {
	l.HandleEvent(discordConnectionErrorMessage, errorMessage)
}

func (l *StandardLogger) DiscordSessionError(errorMessage error) {
	l.HandleEvent(discordSessionErrorMessage, errorMessage)
}

func (l *StandardLogger) SQLError(errorMessage error) {
	l.HandleEvent(sqlErrorMessage, errorMessage)
}

func (l *StandardLogger) DSNDebug(dsn string) {
	l.HandleEvent(dsnDebugMessage, dsn)
}

func (l *StandardLogger) SQLOpeningDebug(dsn string) {
	l.HandleEvent(sqlOpeningMessage, dsn)
}

func (l *StandardLogger) SQLOpenedDebug(dsn string) {
	l.HandleEvent(sqlOpenedMessage, dsn)
}

// TODO: Log output to text file and then work on SQL operations
func (l *StandardLogger) HandleEvent(event Event, extras ...interface{}) {
	if len(extras) > 0 {
		event.message = fmt.Sprintf(event.message.(string), extras[0])
	}

	l.Log(event.level, event.message)
}
