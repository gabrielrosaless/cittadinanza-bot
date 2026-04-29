# Cittadinanza Bot 🇮🇹🤖

Un bot en Go diseñado para monitorear automáticamente la página web de noticias del Consulado de Italia en Córdoba (Argentina), y notificar por Telegram cuando se publiquen avisos relacionados con la apertura de turnos.

## Características
- 🕒 Intervalo de chequeo configurable (Ej: corre cada 8 horas).
- 💾 Base de datos local en SQLite (`bot.db`) para evitar alertas duplicadas.
- 🎯 Macheo de palabras clave insensible a mayúsculas y acentos (ej: `ciudadanía` == `ciudadania`).
- ⚡ Desarrollado en Go para consumir pocos recursos, y compilable en un único archivo binario sin dependencias externas.
- 🔄 Reintenta automáticamente en caso de errores de red (hasta 3 veces por ciclo).

---

## 1. Configurar métodos de notificación
Podés usar **Telegram**, **Email (Gmail)**, o los dos al mismo tiempo. Solo tenés que completar los datos de lo que quieras usar en el `config.json`.

### Opción A: Telegram
Si querés usar Telegram, necesitás dos cosas: un **Token del Bot** y tu **Chat ID**.
1. Abrí Telegram, buscá a **@BotFather** y mandale el mensaje `/newbot`.
2. Seguí sus instrucciones para darle nombre y username a tu bot.
3. @BotFather te va a mandar un **API Token** *(se ve como `123456789:ABCdefGHI...`)*. Guardalo.
4. Iniciá tu nuevo bot enviándole un mensaje cualquiera (ej. "hola").
5. Abrí esta URL en tu navegador reemplazando `<TU_TOKEN>`: `https://api.telegram.org/bot<TU_TOKEN>/getUpdates`
6. Buscá la parte que dice `"chat":{"id":123456789...`. Ese número es tu **Chat ID**.

### Opción B: Email (usando Gmail SMTP)
Si preferís que te mande un mail desde tu propia cuenta de Gmail:
1. Andá a la configuración de seguridad de tu cuenta de Google (necesitás tener la verificación en 2 pasos activada).
2. Buscá la sección de **Contraseñas de aplicaciones** (App Passwords).
3. Creá una nueva contraseña, ponéle el nombre "Cittadinanza Bot" y Google te dará una contraseña rara de 16 letras (ej: `abcd efgh ijkl mnop`).
4. Essa contraseña es la que vas a poner en `email_app_password`.

---

## 3. Subir a GitHub (Serverless gratis)

El bot está diseñado para vivir "en la nube" usando **GitHub Actions** sin gastar un centavo y sin dejar tu computadora prendida.

1. Creá un repositorio privado nuevo en tu cuenta de GitHub.
2. Empujá este código a ese repositorio:
   ```bash
   git init
   git add .
   git commit -m "Initial commit"
   git branch -M main
   git remote add origin https://github.com/TU_USUARIO/TU_REPO.git
   git push -u origin main
   ```
3. En la página de GitHub de tu repositorio, andá a **Settings -> Secrets and variables -> Actions**.
4. Creá los siguientes secretos tocando "New repository secret". Poneles *exactamente* el mismo nombre en mayúscula, y pegales el valor del dato que obtuviste en el Paso 1:
   - `EMAIL_SENDER` (tu mail)
   - `EMAIL_APP_PASSWORD` (la contraseña amarilla de 16 caracteres)
   - `EMAIL_RECIPIENT` (opcional, tu mail o adonde querés que llegue)
   - *(Y si querés usar Telegram, agregá `TELEGRAM_TOKEN` y `TELEGRAM_CHAT_ID`)*

¡Listo! A partir de ahora, GitHub va a despertar el bot mágicamente **cada 8 horas**, y va a revisar si hay turnos. No tenés que hacer absolutamente nada más.
