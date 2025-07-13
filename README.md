# 12-07-2025

curl -v -X POST http://localhost:8080/task

curl -v -X POST -H "Content-Type: application/json" -d '[
{"url":"www.ya.ru"},
{"url":"www.dlya.ru"},
{"url":"www.Nya.ru"},
{"url":"www.Fya.ru"},
{"url":"www.Gya.ru"}]' http://localhost:8080/task/{id}

curl -v -X POST -H "Content-Type: application/json" -d '[
{"url":"https://res.cloudinary.com/startup-grind/image/upload/c_fill,dpr_2.0,f_auto,g_center,h_1080,q_100,w_1080/v1/gcs/platform-data-goog/events/gopherBig.png"},
{"url":"https://cdn-edge.kwork.ru/pics/t3/42/35730409-66e7451e51fd2.jpg"},
{"url":"https://constitutionrf.ru/constitutionrf.pdf"}]' http://localhost:8080/task/{id}

curl -v GET http://localhost:8080/info/{id}


curl -v -X POST -H "Content-Type: application/json" -d '[
{"url":"https://res.cloudinary.com/startup-grind/image/upload/c_fill,dpr_2.0,f_auto,g_center,h_1080,q_100,w_1080/v1/gcs/platform-data-goog/events/gopherBig.png"},
{"url":"https://cdn-edge.kwork.ru/pics/t3/42/35730409-66e7451e51fd2.jpg"}]' http://localhost:8080/task/{id}

curl -v -X POST -H "Content-Type: application/json" -d '[
{"url":"https://res.cloudinary.com/startup-grind/image/upload/c_fill,dpr_2.0,f_auto,g_center,h_1080,q_100,w_1080/v1/gcs/platform-data-goog/events/gopherBig.png"},
{"url":"https://cdn-edge.kwork.ru/pics/t3/42/35730409-66e7451e51fd2.jpg"}]' http://localhost:8080/task/{id}



