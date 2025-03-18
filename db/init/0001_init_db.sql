CREATE USER metric
    PASSWORD 'password';

CREATE DATABASE metrics
    OWNER 'metric'
    ENCODING 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8';