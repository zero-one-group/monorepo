
# Load Balancer and Reverse Proxy via NGINX

A comprehensive tutorial template load balancer configuration


### Prerequisites

- EC2 module nginx provisioning

- Swarm service already deployed

- Template load balancer from monorepo generated to EC2 nginx

## Nginx Deployment Process

### SSL setup

in template there is ssl script **ssl.sh**, there will be 4 option in there generate ssl self signed, generate ssl domain via standalone(without nginx), generate ssl via webroot(nginx container running), generate ssl wildcard. You may update the script as you need for example need more domain. Check this code carefully:
```
DOMAIN="CHANGEME_APPS_DOMAIN"
WILDCARD_DOMAIN="CHANGEME_WILDCARD_DOMAIN"
```
After successfully generated then check on ssl folder there will be generated **.key** and **.crt** there.

Notes: 
Also let say you want to deploy apps first and then monitoring later on, that means you need to comment script that have correlation with monitoring in this case portainer, grafana and nginx monitoring in **default.conf** and **deploy.sh**. Also if you think **WAF** way to aggresive then there is line code:
```
LearningMode;
```
This means waf in watch mode and not block any potential malicious according to core rules. But this will catch on **logs/error.log** so you can make whitelisting rules based on that **error.log** if it false positive. You can enable **WAF** again with comment line above

### Deploy script

There is bash script called **deploy.sh** you may adjust this as you needs if you want to add more apps or other tools monitoring. If already setup then run this bash script. If you want to deploy nginx monitoring then you have deploy first on **goaccess** folder. after that you may deploy again using **deploy.sh**.




