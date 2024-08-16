FROM nginx:1.27.1

WORKDIR /var/www/html/

COPY /nginx/nginx.conf /etc/nginx/nginx.conf

COPY /templates/ /var/www/html/
COPY /static/ /var/www/html/static/


EXPOSE 80