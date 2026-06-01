# test/vulnerable-app/Dockerfile.app
FROM maven:3.9.1-eclipse-temurin-17-alpine AS build
WORKDIR /app
COPY pom.xml .
COPY src ./src
RUN mvn clean package -DskipTests

FROM eclipse-temurin:17-jre-alpine
WORKDIR /app
COPY --from=build /app/target/vulnerable-app-simple-1.0.0.jar app.jar

# Устанавливаем curl для простейшей проверки связи
RUN apk add --no-url-cache curl || apk add --no-cache curl

# Записываем секретный флаг в файловую систему ОС контейнера
RUN mkdir -p /var/lib/secret && echo "FLAG{REAL_LOG4SHELL_FIND_RCE_CONTAINER_2026}" > /var/lib/secret/flag.txt

EXPOSE 8080
ENTRYPOINT ["java", "-Dcom.sun.jndi.ldap.object.trustURLCodebase=true", "-jar", "app.jar"]
