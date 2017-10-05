FROM ubuntu:16.04

RUN apt-get update && apt-get -y install curl
RUN apt-get -y install nodejs npm
RUN ln -s /usr/bin/nodejs /usr/bin/node
RUN npm install -g n
RUN n stable
# Test framework
RUN npm install -g jasmine-node@1.14.5
RUN npm install -g frisby@0.8.5
# Test utilities
RUN npm install -g mysql@2.14.1
RUN npm install -g ssh2-sftp-client@1.1.0
RUN npm install -g http-status-codes@1.2.0
RUN npm install -g form-data@2.2.0
RUN npm install -g path
RUN npm install -g fs
RUN npm install -g cf-nodejs-client@0.13.0

# Configuration for Copy repository
RUN apt-get update && apt-get install -y openssh-server
RUN mkdir /var/run/sshd
RUN echo 'root:screencast' | chpasswd
RUN sed -i 's/PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config

# SSH login fix. Otherwise user is kicked off after login
RUN sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd

# Open SSH port
ENV NOTVISIBLE "in users profile"
RUN echo "export VISIBLE=now" >> /etc/profile
EXPOSE 22
CMD ["/usr/sbin/sshd", "-D"]

# Add non root user
RUN useradd -d /home/frisby -m frisby
RUN echo 'frisby:test' | chpasswd

# Folder where the copies are stored
RUN mkdir /var/lib/crb
RUN chmod a-rw /var/lib/crb

# Load the existing JSfunctional
RUN mkdir /root/functionaltests
WORKDIR /root/functionaltests
ADD * /root/functionaltests/
