FROM python:3

WORKDIR /app

RUN apt-get install -y wget
# Chrome dependency Instalation
RUN apt-get update && apt-get install -y \
  fonts-liberation \
  libasound2 \
  libatk-bridge2.0-0 \
  libatk1.0-0 \
  libatspi2.0-0 \
  libcups2 \
  libdbus-1-3 \
  libdrm2 \
  libgbm1 \
  libgtk-3-0 \
  #    libgtk-4-1 \
  libnspr4 \
  libnss3 \
  libwayland-client0 \
  libxcomposite1 \
  libxdamage1 \
  libxfixes3 \
  libxkbcommon0 \
  libxrandr2 \
  xdg-utils \
  libu2f-udev \
  libvulkan1
RUN wget -q https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb

# Google Chrome
# RUN apt-get install -y ./google-chrome-stable_current_amd64.deb
# RUN rm google-chrome-stable_current_amd64.deb
# Mysql stuff
# RUN apt-get -y install default-libmysqlclient-dev build-essential pkg-config
RUN apt-get install -y python3-dev default-libmysqlclient-dev build-essential

COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt
RUN pip install webdriver-manager
# RUN pip install mysqlclient
# RUN pip install requests[security]

COPY . .
RUN chmod +x ./chromedriver

CMD [ "python", "./main.py" ]