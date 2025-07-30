package system

import (
	"api-server/internal/app/app/system/announcement"
	"api-server/internal/app/app/system/department"
	"api-server/internal/app/app/system/menu"
	"api-server/internal/app/app/system/nav"
	"api-server/internal/app/app/system/role"
	"api-server/internal/app/app/system/user"
	"api-server/internal/middleware/middleware"

	"github.com/gin-gonic/gin"
)

// InitSystemRoutes 初始化系统模块路由
func InitSystemRoutes(r *gin.RouterGroup) {
	// 用户认证相关路由（无需JWT验证，但需要限流保护）
	authGroup := r.Group("/auth")
	{
		// 登录接口使用专门的登录限流
		authGroup.POST("/login", middleware.LoginRateLimit(), user.Login)
		// 验证码接口使用一般限流（支持GET和POST）
		authGroup.GET("/captcha", middleware.RateLimit(10, 20), user.GetCaptcha)
		authGroup.POST("/captcha", middleware.RateLimit(10, 20), user.GetCaptcha)
		// 登出接口使用一般限流
		authGroup.POST("/logout", middleware.RateLimit(5, 10), user.Logout)
	}

	// 需要JWT验证的路由 + 一般限流保护
	protectedGroup := r.Group("/")
	protectedGroup.Use(middleware.JWTAuth())
	protectedGroup.Use(middleware.RateLimit(100, 200)) // 每秒100次请求，突发200次
	{
		// 用户管理
		userGroup := protectedGroup.Group("/users")
		{
			userGroup.GET("", user.GetUserListAPI)            // 获取用户列表
			userGroup.POST("", user.CreateUserAPI)            // 创建用户
			userGroup.GET("/info", user.GetUserInfo)          // 获取当前用户信息
			userGroup.PUT("/info", user.UpdateUserInfo)       // 更新当前用户信息
			userGroup.GET("/menu", user.GetUserMenuList)      // 获取当前用户菜单
			userGroup.POST("/password", user.ChangePassword)  // 修改密码
			userGroup.POST("/avatar", user.UploadAvatarAPI)   // 上传头像
			userGroup.DELETE("/avatar", user.DeleteAvatarAPI) // 删除头像
			userGroup.PUT("/:id", user.UpdateUserAPI)         // 更新用户
			userGroup.DELETE("/:id", user.DeleteUserAPI)      // 删除用户
			userGroup.GET("/:id", user.GetUserDetail)         // 获取用户详情
		}

		// 角色管理
		roleGroup := protectedGroup.Group("/roles")
		{
			roleGroup.GET("", role.GetRoleList)       // 获取角色列表
			roleGroup.POST("", role.AddRole)          // 创建角色
			roleGroup.PUT("", role.UpdateRole)        // 更新角色
			roleGroup.DELETE("", role.DeleteRole)     // 删除角色
			roleGroup.GET("/:id", role.GetRoleDetail) // 获取角色详情
		}

		// 菜单管理
		menuGroup := protectedGroup.Group("/menus")
		{
			menuGroup.GET("", menu.GetMenuList)                 // 获取菜单列表
			menuGroup.POST("", menu.AddMenu)                    // 创建菜单
			menuGroup.PUT("", menu.UpdateMenu)                  // 更新菜单
			menuGroup.DELETE("", menu.DeleteMenu)               // 删除菜单
			menuGroup.GET("/:id", menu.GetMenuDetail)           // 获取菜单详情
			menuGroup.GET("/tree", menu.GetMenuTree)            // 获取菜单树
			menuGroup.GET("/role", menu.GetMenuListByRoleID)    // 获取角色菜单权限
			menuGroup.PUT("/role", menu.UpdateMenuListByRoleID) // 更新角色菜单权限
		}

		// 菜单权限管理
		menuAuthGroup := protectedGroup.Group("/menu-auths")
		{
			menuAuthGroup.GET("", menu.GetMenuAuthList)   // 获取菜单权限列表
			menuAuthGroup.POST("", menu.AddMenuAuth)      // 添加菜单权限
			menuAuthGroup.PUT("", menu.UpdateMenuAuth)    // 更新菜单权限
			menuAuthGroup.DELETE("", menu.DeleteMenuAuth) // 删除菜单权限
		}

		// 部门管理
		deptGroup := protectedGroup.Group("/departments")
		{
			deptGroup.GET("", department.GetDepartmentList)       // 获取部门列表
			deptGroup.POST("", department.AddDepartment)          // 创建部门
			deptGroup.PUT("", department.UpdateDepartment)        // 更新部门
			deptGroup.DELETE("", department.DeleteDepartment)     // 删除部门
			deptGroup.GET("/:id", department.GetDepartmentDetail) // 获取部门详情
		}

		// 导航管理
		navGroup := protectedGroup.Group("/navs")
		{
			navGroup.GET("", nav.GetNavListAPI)                   // 获取导航列表
			navGroup.GET("/:id", nav.GetNavDetail)                // 获取导航详情
			navGroup.POST("", nav.CreateNavAPI)                   // 创建导航
			navGroup.PUT("/:id", nav.UpdateNavAPI)                // 更新导航
			navGroup.DELETE("/:id", nav.DeleteNavAPI)             // 删除导航
			navGroup.POST("/sort", nav.SortNavAPI)                // 导航排序
			navGroup.GET("/user/:user_id", nav.GetUserNavListAPI) // 获取用户导航
			navGroup.GET("/role/:role_id", nav.GetRoleNavListAPI) // 获取角色导航
			navGroup.GET("/public", nav.GetPublicNavListAPI)      // 获取公共导航
		}

		// 公告管理
		announcementGroup := protectedGroup.Group("/announcements")
		{
			announcementGroup.GET("", announcement.GetAnnouncementListAPI)           // 获取公告列表
			announcementGroup.POST("", announcement.CreateAnnouncementAPI)           // 创建公告
			announcementGroup.GET("/active", announcement.GetActiveAnnouncementsAPI) // 获取活跃公告
			announcementGroup.GET("/:id", announcement.GetAnnouncementDetailAPI)     // 获取公告详情
			announcementGroup.PUT("/:id", announcement.UpdateAnnouncementAPI)        // 更新公告
			announcementGroup.DELETE("/:id", announcement.DeleteAnnouncementAPI)     // 删除公告
		}
	}
}
