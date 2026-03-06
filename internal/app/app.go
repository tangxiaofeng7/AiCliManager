package app

import (
	"AiCliManager/internal/service"
	"context"
	"database/sql"
)

// App 是 Wails 绑定的核心结构体，所有前端可调用的方法都挂载在此
type App struct {
	ctx               context.Context
	db                *sql.DB
	providerService   *service.ProviderService
	profileService    *service.ProfileService
	proxyService      *service.ProxyService
	mcpService        *service.McpService
	skillService      *service.SkillService
	sessionService    *service.SessionService
	cliSessionService *service.CliSessionService
	settingsService   *service.SettingsService
	terminalService   *service.TerminalService
	syncService       *service.SyncService
	launcherService   *service.LauncherService
}

// New 创建并初始化 App 实例，依赖注入所有 Service
func New(db *sql.DB) *App {
	providerSvc := service.NewProviderService(db)
	proxySvc := service.NewProxyService(db)
	mcpSvc := service.NewMcpService(db)
	skillSvc := service.NewSkillService(db)
	sessionSvc := service.NewSessionService(db)
	cliSessionSvc := service.NewCliSessionService()
	terminalSvc := service.NewTerminalService()
	syncSvc := service.NewSyncService()
	profileSvc := service.NewProfileService(db)
	settingsSvc := service.NewSettingsService(db)

	launcherSvc := service.NewLauncherService(
		db,
		providerSvc,
		proxySvc,
		mcpSvc,
		skillSvc,
		sessionSvc,
		terminalSvc,
		syncSvc,
	)

	return &App{
		db:                db,
		providerService:   providerSvc,
		profileService:    profileSvc,
		proxyService:      proxySvc,
		mcpService:        mcpSvc,
		skillService:      skillSvc,
		sessionService:    sessionSvc,
		cliSessionService: cliSessionSvc,
		settingsService:   settingsSvc,
		terminalService:   terminalSvc,
		syncService:       syncSvc,
		launcherService:   launcherSvc,
	}
}

// Startup 在 Wails 应用启动时调用，保存 context
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// Shutdown 在 Wails 应用关闭时调用，释放资源
func (a *App) Shutdown(ctx context.Context) {
	if a.db != nil {
		_ = a.db.Close()
	}
}
