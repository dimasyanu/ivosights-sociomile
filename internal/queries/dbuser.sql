CREATE USER 'admin'@'%' IDENTIFIED BY 'admin_password';
GRANT ALL PRIVILEGES ON *.* TO 'admin'@'%';
FLUSH PRIVILEGES;
