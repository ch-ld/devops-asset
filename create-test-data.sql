-- 创建测试主机组数据
INSERT INTO cmdb_host_groups (id, name, description, parent_id, path, sort, created_at, updated_at) VALUES
(1, '生产环境', '生产环境主机组', NULL, '/生产环境', 1, NOW(), NOW()),
(2, 'Web服务器', 'Web服务器组', 1, '/生产环境/Web服务器', 1, NOW(), NOW()),
(3, '数据库服务器', '数据库服务器组', 1, '/生产环境/数据库服务器', 2, NOW(), NOW()),
(4, '测试环境', '测试环境主机组', NULL, '/测试环境', 2, NOW(), NOW()),
(5, '开发环境', '开发环境主机组', NULL, '/开发环境', 3, NOW(), NOW());

-- 查看创建的数据
SELECT * FROM cmdb_host_groups ORDER BY sort, id;
