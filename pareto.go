package pareto

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"github.com/zourva/pareto/box/env"
	"github.com/zourva/pareto/box/prof"
	"github.com/zourva/pareto/logger"
	"os"
)

type paretoKit struct {
	workingDir *env.WorkingDir
	logger     *logger.Logger
	//diagnoser  *diagnoser.Diagnoser
	//monitor    *monitor.SysMonitor
	//updater    *updater.OtaManager
	profiler  *prof.Profiler
	flagParse bool
}

var bot = new(paretoKit)

// Option defines pareto initialization options.
type Option func()

// EnableFlagParse enables or disables flag.Parse.
func EnableFlagParse(parse bool) Option {
	return func() {
		bot.flagParse = parse
		if parse {
			flag.Parse()
		}
	}
}

// EnableProfiler enables prof.Profiler.
func EnableProfiler() Option {
	return func() {
		bot.profiler = prof.NewProfiler(nil)
		bot.profiler.Start()
	}
}

// WithLogger allows to provide a logger config
// as an option.
func WithLogger(l *logger.Logger) Option {
	return func() {
		bot.logger = l
	}
}

// WithWorkingDirLayout allows to hint working dir layout.
func WithWorkingDirLayout(wd *env.WorkingDir) Option {
	return func() {
		bot.workingDir = wd
	}
}

// WithWorkingDir acts the same as WithWorkingDirLayout except that
// it also set system level working dir, using os.Chdir, to the parent
// directory of this executable.
func WithWorkingDir(wd *env.WorkingDir) Option {
	return func() {
		bot.workingDir = wd
		err := os.Chdir(env.GetExecFilePath() + "/../")
		if err != nil {
			log.Fatalln("change working dir failed:", err)
		}
	}
}

//// WithDiagnoser allows to provide a diagnoser service config
//// as an option.
//func WithDiagnoser(d *diagnoser.Diagnoser) Option {
//	return func() {
//		bot.diagnoser = d
//	}
//}

//// WithUpdater allows to provide an updater service config
//// as an option.
//func WithUpdater(u *updater.OtaManager) Option {
//	return func() {
//		bot.updater = u
//	}
//}

//// WithMonitor allows to provide a monitor service config
//// as an option.
//func WithMonitor(m *monitor.SysMonitor) Option {
//	return func() {
//		bot.monitor = m
//	}
//}

// WithCli allows to provide a command line interface component config
// as an option.
func WithCli() Option {
	return func() {
	}
}

// SetupWithOpts create a pareto environment with the
// given options.
func SetupWithOpts(options ...Option) {
	for _, o := range options {
		o()
	}

	log.Infoln("setup pareto environment done")
}

//// Setup creates a default logger and working dir,
//// enables flag.Parse
//func Setup() {
//	SetupWithOpts(
//		EnableFlagParse(true),
//		WithLogger(
//			logger.NewLogger(&logger.Options{
//				Verbosity:   "v",
//				LogFileName: env.GetExecFilePath() + "/../log/out.log",
//				MaxSize:     50,
//				MaxAge:      7,
//				MaxBackups:  3,
//			}),
//		),
//		WithWorkingDir(
//			env.NewWorkingDir(true,
//				[]*env.DirInfo{
//					{Name: "bin", Mode: 0755},
//					{Name: "etc", Mode: 0755},
//					{Name: "lib", Mode: 0755},
//					{Name: "log", Mode: 0755},
//				}),
//		))
//}

// Teardown tears down the working space
func Teardown() {
	if bot.profiler != nil {
		bot.profiler.Stop()
	}

	log.Infoln("teardown pareto environment done")
}
