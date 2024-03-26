### Relations

#### User
- **Id**: Уникальный идентификатор пользователя в базе данных.
- **Login**: Электронная почта пользователя, используемая для входа.
- **Password**: Пароль пользователя.
- **Name**: Имя пользователя.
- **Surname**: Фамилия пользователя.
- **Middlename**: Отчество пользователя.
- **Gender**: Пол пользователя.
- **Birthday**: Дата рождения пользователя.
- **RegistrationDate**: Дата регистрации пользователя.
- **AvatarId**: Ссылка на фотографию пользователя.
- **PhoneNumber**: Номер телефона пользователя.
- **Description**: Дополнительная информация, которую пользователь может предоставить о себе.

#### Email
- **Id**: Уникальный идентификатор письма в базе данных.
- **Topic**: Тема письма.
- **Text**: Текст письма.
- **DateOfDispatch**: Дата отправки письма.
- **PhotoId**: Ссылка на фотографию отправителя письма.
- **SenderId**: Уникальный идентификатор пользователя, отправившего письмо.
- **RecipientId**: Уникальный идентификатор пользователя, получившего письмо.
- **ReadStatus**: Статус прочтения письма (прочтено/непрочтено).
- **DeletedStatus**: Статус удаления письма (в корзине/не в корзине).
- **DraftStatus**: Статус черновика письма (черновик/не черновик).
- **ReplyToEmailId**: Уникальный идентификатор письма, на который данное письмо является ответом (если есть).
- **Flag**: Флаг, который может быть установлен пользователем (например, помечено как важное).

#### File
- **Id**: Уникальный идентификатор вложения в базе данных.
- **EmailId**: Уникальный идентификатор письма, к которому прикреплено вложение.
- **DocumentId**: Ссылка на документ.
- **VideoId**: Ссылка на видео.
- **GifId**: Ссылка на гифку.
- **MusicId**: Ссылка на музыку.
- **ArchiveId**: Ссылка на архив.

#### UserEmail
- **Id**: Уникальный идентификатор связи пользователя с письмом в базе данных.
- **UserId**: Уникальный идентификатор пользователя, участвующего в переписке.
- **EmailId**: Уникальный идентификатор письма, полученного или отправленного пользователем.

#### Folder
- **Id**: Уникальный идентификатор папки в базе данных.
- **UserId**: Уникальный идентификатор пользователя, которому принадлежит папка.
- **Name**: Название папки.

#### FolderEmail
- **Id**: Уникальный идентификатор связи папки с письмом в базе данных.
- **FolderId**: Уникальный идентификатор папки, в которой находится письмо.
- **EmailId**: Уникальный идентификатор письма, находящегося в папке.

#### Settings
- **Id**: Уникальный идентификатор настроек пользователя в базе данных.
- **UserId**: Уникальный идентификатор пользователя, которому принадлежат настройки.
- **NotificationTolerance**: Статус уведомлений пользователя (включены/выключены).
- **Language**: Язык интерфейса пользователя.

#### Session
- **Id**: Уникальный идентификатор сессии пользователя в базе данных.
- **UserId**: Уникальный идентификатор пользователя, которому принадлежит сессия.
- **CreationDate**: Дата и время создания сессии.
- **Device**: Устройство, с которого была инициирована сессия.
- **LifeTime**: Время действия сессии.
- **CsrfToken**: Токен CSRF, используемый для защиты от атак межсайтовой подделки запросов.

---
Simple ER-diagram
---

```mermaid
erDiagram
    USER ||--o{ SESSION : "Owns"
    USER ||--o{ SETTINGS : "Has"
    USER ||--o{ FOLDER : "Owns"
    FOLDER ||--o{ FOLDEREMAIL : "Contains"
    EMAIL ||--o{ FOLDEREMAIL : "Located"
    EMAIL ||--o{ USEREMAIL : "Related"
    USER ||--o{ USEREMAIL : "Received"
    EMAIL ||--|{ FILE : "Contains"
```

---
ER-diagram
---

```mermaid
erDiagram
    USER {
        string Id
        string Login
        string Password
        string Name
        string Surname
        string Middlename
        string Gender
        date Birthday
        datetime RegistrationDate
        string AvatarId
        string PhoneNumber
        string Description
    }
    EMAIL {
        string Id
        string Topic
        string Text
        datetime DateOfDispatch
        string PhotoId
        string SenderId
        string RecipientId
        boolean ReadStatus
        boolean DeletedStatus
        boolean DraftStatus
        string ReplyToEmailId
        string Flag
    }
    FILE {
        string Id
        string EmailId
        string DocumentId
        string VideoId
        string GifId
        string MusicId
        string ArchiveId
    }
    USEREMAIL {
        string Id
        string UserId
        string EmailId
    }
    FOLDER {
        string Id
        string UserId
        string Name
    }
    FOLDEREMAIL {
        string Id
        string FolderId
        string EmailId
    }
    SETTINGS {
        string Id
        string UserId
        boolean NotificationTolerance
        string Language
    }
    SESSION {
        string Id
        string UserId
        datetime CreationDate
        string Device
        duration LifeTime
        string CsrfToken
    }

    USER ||--o{ SESSION : "Owns"
    USER ||--o{ SETTINGS : "Has"
    USER ||--o{ FOLDER : "Owns"
    FOLDER ||--o{ FOLDEREMAIL : "Contains"
    EMAIL ||--o{ FOLDEREMAIL : "Located"
    EMAIL ||--o{ USEREMAIL : "Related"
    USER ||--o{ USEREMAIL : "Received"
    EMAIL ||--|{ FILE : "Contains"
```

```mermaid
erDiagram
    USER {
        Id
        Login
        Password
        Name
        Surname
        Middlename
        Gender
        Birthday
        RegistrationDate
        AvatarId
        PhoneNumber
        Description
    }
    EMAIL {
        Id
        Topic
        Text
        DateOfDispatch
        PhotoId
        SenderId
        RecipientId
        ReadStatus
        DeletedStatus
        DraftStatus
        ReplyToEmailId
        Flag
    }
    FILE {
        Id
        EmailId
        DocumentId
        VideoId
        GifId
        MusicId
        ArchiveId
        _
    }
    USEREMAIL {
        Id
        UserId
        EmailId
        _
    }
    FOLDER {
        Id
        UserId
        Name
        _
    }
    FOLDEREMAIL {
        Id
        FolderId
        EmailId
        _
    }
    SETTINGS {
        Id
        UserId
        NotificationTolerance
        Language
    }
    SESSION {
        Id
        UserId
        CreationDate
        Device
        LifeTime
        CsrfToken
    }

    USER ||--o{ SESSION : "Owns"
    USER ||--o{ SETTINGS : "Has"
    USER ||--o{ FOLDER : "Owns"
    FOLDER ||--o{ FOLDEREMAIL : "Contains"
    EMAIL ||--o{ FOLDEREMAIL : "Located"
    EMAIL ||--o{ USEREMAIL : "Related"
    USER ||--o{ USEREMAIL : "Received"
    EMAIL ||--|{ FILE : "Contains"
```