package app

// GetSetting 读取单条全局设置，不存在时返回 defaultValue
func (a *App) GetSetting(key, defaultValue string) (string, error) {
	return a.settingsService.Get(key, defaultValue)
}

// SetSetting 写入或更新一条全局设置
func (a *App) SetSetting(key, value string) error {
	if key == "" {
		return errInvalidParam("设置 key 不能为空")
	}
	return a.settingsService.Set(key, value)
}

// GetAllSettings 获取所有设置项（key-value map）
func (a *App) GetAllSettings() (map[string]string, error) {
	return a.settingsService.GetAll()
}

// DeleteSetting 删除一条设置
func (a *App) DeleteSetting(key string) error {
	if key == "" {
		return errInvalidParam("设置 key 不能为空")
	}
	return a.settingsService.Delete(key)
}
