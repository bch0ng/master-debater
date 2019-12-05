# add the necessary instructions
# to create a Docker container image
# for your summary client using the
# NGINX web server image as the base
FROM nginx
#copies this directories contents into nginx's html folder
COPY . /usr/share/nginx/html/
ADD default.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
EXPOSE 443
EXPOSE 3000
#chmod +x /usr/share/nginx/html
#ENTRYPOINT ["/usr/share/nginx/html"]
#
