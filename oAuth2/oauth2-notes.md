# Learning process
* go installed
* vagrant installed
* git repo downloaded (https://github.com/Charnnarong/oauthsandbox) under /Users/vytautas/dev/Udemy/oAuth2/oauthsandbox
* vagrant installed using
* configure network in VirtualBox:
  - Adapter 1 - NAT
  - Adapter 2 - Host-only Adapter, name: vboxnet0 (or that which corresponds IP address in Vagrant file)
```vagrant up```
* keycload downloaded from official site and installed inside virtualBox in the folder _/home/vagrant/keycloak-11-0.2_
* configure IP address of Keycloak under _./standalone/configuration/standalone.xml_, address
* start Keycloak using command:
```
./bin/standalone.sh
```

VM username:password = vagrant:vagrant

# Keycloak info:
* Address http://192.168.56.10:8080/
* Admin username:password = admin:admin

# Start client (golang)
Address: http://localhost:8081



# Vagrant commands:
```vagrant up```

```vagrant status```

```vagrant ssh``` - get inside box

```vagrant suspend```

```vagrant resume```

# Unix commands
```
$ netstat -ntlp

$ ip route show
```