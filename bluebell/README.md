# bluebell项目后端

网易云课堂

https://study.163.com/course/courseMain.htm?courseId=1210171207


[docker方式]
需要环境

docker run --name mysql0507 -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root -v /Users/q1mi/docker/mysql:/var/lib/mysql -d mysql:5.7

docker run --name redis0500 -p 6379:6379 -d redis:5.0.14-alpine3.15

编译并运行

docker build . -t bluebell

docker run -p 8888:8888 --link=mysql0507:mysql0507 --link=redis0500:redis0500  bluebell