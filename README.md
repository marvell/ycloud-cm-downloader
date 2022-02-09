# ycloud-cm-downloader

Используется для выгрузки сертификатов из Яндекс.Облака Certificate Manager.

- [Как создать сервисный аккаунт?](https://cloud.yandex.ru/docs/iam/concepts/users/service-accounts)
- [Как создать авторизованный ключ?](https://cloud.yandex.ru/docs/iam/operations/authorized-key/create)

```bash
$ ycloud-cm-downloader -cert-id <cert-id> -sa-id <sa-id> -key-id <key-id> -privkey <privkey-path>
full chain file was created in ./<cert-id>_fullchain.pem
private key file was created in ./<cert-id>_privkey.pem
```