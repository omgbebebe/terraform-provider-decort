/*
Пример использования
Работы с ресурсом basic service group
Ресурс позволяет:
1. Создавать группы
2. Редактировать группы
3. Удалять группы
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

resource "decort_bservice_group" "bsg" {
  #id back service
  #обязательный параметр
  #тип - число
  service_id = 444444

  #название группы
  #обязательный параметр
  #тип - строка
  compgroup_name = "tf_group_rename"

  #id группы
  #необязательный параметр
  #тип - число
  #применяется при редактировании группы, либо при создании .tfstate - файла, если группа имеется в плафторме
  compgroup_id = 33333

  #кол-во вычислительных ресурсов
  #обязательный параметр
  #тип - число
  #используется так же для редактирования группы
  comp_count = 1

  #кол-во ядер на выч. ресурс
  #обязательный параметр
  #тип - число
  #используется так же для редактирования группы
  cpu = 2

  #кол-во оперативной памяти на выч. ресурс, в МБ
  #обязательный параметр
  #тип - число
  #используется так же для редактирования группы
  ram = 256

  #размер диска для выч. ресурса, в ГБ
  #обязательный параметр
  #тип - число
  #используется так же для редактирования группы
  disk = 11

  #id образа диска
  #обязательный параметр
  #тип - число
  image_id = 2222

  #драйвер
  #обязательный параметр
  #тип - число
  driver = "kvm_x86"

  #id сетей extnet
  #обязательный параметр
  #тип - массив чисел
  #должен быть использован vins или extnets
  extnets = [1111]

  #id сетей vinses
  #обязательный параметр
  #тип - массив чисел
  #должен быть использован vins или extnets
  #vinses        = [1111, 2222]

  #время таймуата перед стартом
  #необязательный параметр
  #тип - число
  #используется при создании ресурса
  #timeout_start  = 0

  #тег группы
  #необязательный параметр
  #тип - строка
  #используется при создании и редактировании ресурса
  # role           = "tf_test_changed"

  #id групп родителей
  #необязательный параметр
  #тип - массив чисел
  #используется при редактировании ресурса
  #parents = []

  #принудительное обновление параметров выч. мощностей  (ram,disk,cpu) и имени группы
  #необязательный параметр
  #тип - булев тип
  #используется при редактировании
  #force_update   = true

  #старт/стоп вычислительных мощностей
  #необязательный параметр
  #тип - булев тип
  #используется при редактировании
  #по-умолчанию - false
  #start = false

  #принудительная остановка вычислительных мощностей
  #необязательный параметр
  #тип - булев тип
  #используется при редактировании и остановке группы
  #по-умолчанию - false
  #force_stop = false

  #удаление вычислительных мощностей
  #необязательный параметр
  #тип - массив чисел
  #используется при редактировании
  #remove_computes = [32287]

  #режим увеличения числа выч. мощностей
  #необязательный параметр
  #тип - число
  #используется в связке с comp_count при редактировании группы
  #возможные варианты - RELATIVE и ABSOLUTE
  #mode = "RELATIVE"

}

output "test" {
  value = decort_bservice_group.bsg
}
