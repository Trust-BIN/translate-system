package hanlders

// 注销账号处理函数
//func DeleteAccountHandler(c *gin.Context) {
//	if c.Request.Method != http.MethodPost {
//		c.JSON(http.StatusMethodNotAllowed, serve.ErrorResponse{Error: "只支持 POST 请求"})
//		return
//	}
//
//	username, err := serve.GetUsernameFromCookie(c)
//	if err != nil {
//		return
//	}
//	if username == "" {
//		c.JSON(http.StatusUnauthorized, serve.ErrorResponse{Error: "用户未登录"})
//		return
//	}
//
//	db := database.GetDB()
//
//	// 使用事务确保数据一致性
//	tx, err := db.Begin()
//	if err != nil {
//		log.Printf("开启事务失败: %v", err)
//		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
//		return
//	}
//	defer tx.Rollback() // 确保在出错时回滚
//
//	// 1. 删除用户的历史记录
//	_, err = tx.Exec("DELETE FROM translation_history WHERE username = ?", username)
//	if err != nil {
//		log.Printf("删除用户历史记录时出错: %v", err)
//		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
//		return
//	}
//
//	// 2. 删除用户信息
//	_, err = tx.Exec("DELETE FROM users WHERE username = ?", username)
//	if err != nil {
//		log.Printf("删除用户信息时出错: %v", err)
//		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
//		return
//	}
//
//	// 提交事务
//	if err = tx.Commit(); err != nil {
//		log.Printf("提交事务时出错: %v", err)
//		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
//		return
//	}
//
//	// 3. 清除用户的会话信息
//	c.SetCookie("username", "", -1, "/", "", false, true)
//
//	c.JSON(http.StatusOK, gin.H{
//		"success": true,
//		"message": "账号注销成功",
//	})
//}
