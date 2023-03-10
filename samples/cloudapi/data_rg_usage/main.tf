/*
Пример использования
Получение информации об использовании ресурсов на ресурсной группе
*/
#Расскомментируйте этот код,
#и внесите необходимые правки в версию и путь,
#чтобы работать с установленным вручную (не через hashicorp provider registry) провайдером

/* terraform {
  required_providers {
    decort = {
      version = "1.1"
      source  = "digitalenergy.online/decort/decort"
    }
  }
} */

provider "decort" {
  authenticator = "oauth2"
  #controller_url = <DECORT_CONTROLLER_URL>
  controller_url = "https://ds1.digitalenergy.online"
  #oauth2_url = <DECORT_SSO_URL>
  oauth2_url           = "https://sso.digitalenergy.online"
  allow_unverified_ssl = true
}

data "decort_rg_usage" "rg_usage" {
  #id ресурсной группы
  #обязательный параметр
  #тип - число
  rg_id = 123
  #причина для действия
  #необязательный параметр
  #тип - строка
  reason = "TEST"
}

output "output" {
  value = data.decort_rg_usage.rg_usage
}
