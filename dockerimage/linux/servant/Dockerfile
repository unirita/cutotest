FROM centos
MAINTAINER hondata <takahiro_honda@unirita.co.jp>
ADD cuto /cuto/
RUN mkdir /cuto/joblog
RUN mkdir /cuto/log
ENV CUTOROOT /cuto
EXPOSE 2015
CMD /cuto/bin/servant
