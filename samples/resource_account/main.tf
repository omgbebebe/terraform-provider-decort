/*
Пример использования
Ресурса account
Ресурс позволяет:
1. Создавать аккаунт
2. Редактировать аккаунт
3. Удалять аккаунт

*/
#Расскомментируйте этот код,
#и внесите необходимые правки в версию и путь,
#чтобы работать с установленным вручную (не через hashicorp provider registry) провайдером
/*
terraform {
  required_providers {
    decort = {
      version = "1.1"
      source  = "digitalenergy.online/decort/decort"
    }
  }
}
*/

provider "decort" {
  authenticator = "oauth2"
  #controller_url = <DECORT_CONTROLLER_URL>
  controller_url = "https://ds1.digitalenergy.online"
  #oauth2_url = <DECORT_SSO_URL>
  oauth2_url           = "https://sso.digitalenergy.online"
  allow_unverified_ssl = true
}

resource "decort_account" "a" {
  #имя аккаунта
  #обязательный параметр
  #тип - строка
  #используется при создании и редактировании аккаунта
  account_name = "new_my_account"

  #имя пользователя - создателя аккаунта
  #обязательный параметр
  #тип - строка
  username = "username@decs3o"

  #доступность аккаунта
  #необязательный параметр
  #тип - булев тип
  #может применяться при редактировании аккаунта
  enable = true

  #id аккаунта, позволяет сформировать .tfstate, если аккаунт имеет в платформе
  #необязательный параметр
  #тип - число
  account_id = 11111

  #электронная почта, на которую будет отправлена информация о доступе
  #необязательный параметр
  #тип - строка
  #применяется при создании аккаунта
  emailaddress = "fff@fff.ff"

  #отправлять ли на электронную почту письмо о доступе
  #необязательный параметр
  #тип - булев тип
  #применяется при создании аккаунта и редактировании аккаунта
  send_access_emails = true

  #добавление/редактирование/удаление пользователей, к которым привязан аккаунт
  #необязательный параметр
  #тип - объект, кол-во таких объектов не ограничено
  /*users {
    #id пользователя
    #обязательный параметр
    #тип - строка
    user_id     = "username_2@decs3o"

    #тип доступа пользователя
    #обязательный параметр
    #тип - строка
    #возможные параметры:
    #R - чтение
    #RCX - запись
    #ARCXDU - админ
    access_type = "R"

    #рекурсивное удаление пользователя из всех ресурсов аккаунтов
    #необязательный параметр
    #тип - булев тип
    #по-умолчанию - false
    #применяется при удалении пользователя из аккаунта
    recursive_delete = true
  }
  users {
    user_id     = "username_1@decs3o"
    access_type = "R"
  }*/

  #ограничение используемых ресурсов
  #необязательный параметр
  #тип - объект
  #используется при создании и редактировании
  resource_limits {
    #кол-во используемых ядер cpu
    #необязательный параметр
    #тип - число
    #если установлена -1 - кол-во неограиченно
    cu_c = 2

    #кол-во используемой RAM в МБ
    #необязательный параметр
    #тип - число
    #если установлена -1 - кол-во неограиченно
    cu_m = 1024

    #размер дисков, в ГБ
    #необязательный параметр
    #тип - число
    #если установлена -1 - размер неограичен
    cu_d = 23

    #кол-во используемых публичных IP
    #необязательный параметр
    #тип - число
    #если установлена -1 - кол-во неограиченно
    cu_i = 2

    #ограничения на кол-во передачи данных, в ГБ
    #необязательный параметр
    #тип - число
    #если установлена -1 - кол-во неограиченно
    cu_np = 2

    #кол-во графических процессоров
    #необязательный параметр
    #тип - число
    #если установлена -1 - кол-во неограиченно
    gpu_units = 2
  }

  #восстановление аккаунта
  #необязательный параметр
  #тип - булев тип
  #применяется к удаленным аккаунтам
  #по-умолчанию - false
  #restore = false

  #мгновеное удаление аккаунта, если да - то аккаунт невозможно будет восстановить
  #необязательный параметр
  #тип - булев тип
  #используется при удалении аккаунта
  #по-умолчанию - false
  #permanently = true
}

output "test" {
  value = decort_account.a
}
