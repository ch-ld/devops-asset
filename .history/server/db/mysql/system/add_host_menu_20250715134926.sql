-- 添加主机管理模块菜单数据

-- 获取当前最大ID
SET @max_id = (SELECT COALESCE(MAX(id), 0) FROM system_menus);

-- 添加CMDB父级菜单（如果不存在）
INSERT INTO system_menus (id, path, name, component, title, icon, show_badge, is_hide, is_hide_tab, is_iframe, keep_alive, is_first_level, status, level, parent_id, sort, created_at, updated_at)
SELECT @max_id + 1, '/cmdb', 'CMDB', '/index/index', 'CMDB资产管理', '&#xe6a0;', 2, 2, 2, 2, 2, 2, 1, 1, 0, 90, NOW(), NOW()
FROM dual
WHERE NOT EXISTS (
    SELECT 1 FROM system_menus WHERE name = 'CMDB' AND parent_id = 0
);

-- 更新最大ID（如果添加了新菜单）
SET @cmdb_id = (SELECT id FROM system_menus WHERE name = 'CMDB' AND parent_id = 0);

-- 添加主机管理菜单
INSERT INTO system_menus (id, path, name, component, title, icon, show_badge, is_hide, is_hide_tab, is_iframe, keep_alive, is_first_level, status, level, parent_id, sort, created_at, updated_at)
SELECT @max_id + 2, 'hosts', 'HostManagement', '/cmdb/host/index', '主机管理', '&#xe73e;', 2, 2, 2, 2, 1, 2, 1, 2, @cmdb_id, 99, NOW(), NOW()
FROM dual
WHERE NOT EXISTS (
    SELECT 1 FROM system_menus WHERE name = 'HostManagement' AND parent_id = @cmdb_id
);

-- 添加主机详情菜单（隐藏菜单）
INSERT INTO system_menus (id, path, name, component, title, icon, show_badge, is_hide, is_hide_tab, is_iframe, keep_alive, is_first_level, status, level, parent_id, sort, created_at, updated_at)
SELECT @max_id + 3, 'host/:id', 'HostDetail', '/cmdb/host/detail', '主机详情', '&#xe74c;', 2, 1, 2, 2, 2, 2, 1, 2, @cmdb_id, 0, NOW(), NOW()
FROM dual
WHERE NOT EXISTS (
    SELECT 1 FROM system_menus WHERE name = 'HostDetail' AND parent_id = @cmdb_id
);

-- 添加编辑主机菜单（隐藏菜单）
INSERT INTO system_menus (id, path, name, component, title, icon, show_badge, is_hide, is_hide_tab, is_iframe, keep_alive, is_first_level, status, level, parent_id, sort, created_at, updated_at)
SELECT @max_id + 4, 'host-edit/:id', 'HostEdit', '/cmdb/host/edit', '编辑主机', '&#xe66e;', 2, 1, 2, 2, 2, 2, 1, 2, @cmdb_id, 0, NOW(), NOW()
FROM dual
WHERE NOT EXISTS (
    SELECT 1 FROM system_menus WHERE name = 'HostEdit' AND parent_id = @cmdb_id
);

-- 添加创建主机菜单（隐藏菜单）
INSERT INTO system_menus (id, path, name, component, title, icon, show_badge, is_hide, is_hide_tab, is_iframe, keep_alive, is_first_level, status, level, parent_id, sort, created_at, updated_at)
SELECT @max_id + 5, 'host-create', 'HostCreate', '/cmdb/host/edit', '添加主机', '&#xe727;', 2, 1, 2, 2, 2, 2, 1, 2, @cmdb_id, 0, NOW(), NOW()
FROM dual
WHERE NOT EXISTS (
    SELECT 1 FROM system_menus WHERE name = 'HostCreate' AND parent_id = @cmdb_id
);

-- 添加主机仪表盘菜单
INSERT INTO system_menus (id, path, name, component, title, icon, show_badge, is_hide, is_hide_tab, is_iframe, keep_alive, is_first_level, status, level, parent_id, sort, created_at, updated_at)
SELECT @max_id + 6, 'host-dashboard', 'HostDashboard', '/cmdb/host/dashboard', '主机概览', '&#xe64c;', 2, 2, 2, 2, 1, 2, 1, 2, @cmdb_id, 88, NOW(), NOW()
FROM dual
WHERE NOT EXISTS (
    SELECT 1 FROM system_menus WHERE name = 'HostDashboard' AND parent_id = @cmdb_id
);

-- 添加云账号管理菜单
INSERT INTO system_menus (id, path, name, component, title, icon, show_badge, is_hide, is_hide_tab, is_iframe, keep_alive, is_first_level, status, level, parent_id, sort, created_at, updated_at)
SELECT @max_id + 7, 'providers', 'ProviderManagement', '/cmdb/provider/index', '云账号管理', '&#xe67c;', 2, 2, 2, 2, 1, 2, 1, 2, @cmdb_id, 77, NOW(), NOW()
FROM dual
WHERE NOT EXISTS (
    SELECT 1 FROM system_menus WHERE name = 'ProviderManagement' AND parent_id = @cmdb_id
);

-- 给超级管理员角色(id=1)添加主机管理权限
INSERT INTO system_roles__system_menus (system_role_id, system_menu_id)
SELECT 1, id FROM system_menus
WHERE name IN ('CMDB', 'HostManagement', 'HostDetail', 'HostEdit', 'HostCreate', 'HostDashboard', 'ProviderManagement')
AND id NOT IN (SELECT system_menu_id FROM system_roles__system_menus WHERE system_role_id = 1); 
