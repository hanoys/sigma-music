#!/bin/bash

# Убедимся, что директория для логов существует
mkdir -p phout

# Функция для запуска тестов с заданным портом
run_tests() {
  local port=$1  # Порт для теста
  local target="localhost:$port"  # Цель с указанием порта

  for i in {1..10}; do
    # Изменяем конфигурационный файл load.yaml с новым портом и сохраняем в него изменения
    sed -i "s/target: localhost:[0-9]*/target: $target/" load.yaml

    # Выполняем тест
    ./pandora_0.5.32_darwin_arm64 load.yaml

    # Перемещаем результат в отдельный файл с указанием порта и номера прогона
    mv phout.log "phout/phout${port}_${i}.log"
    
    echo "Запуск ${i} на ${target} завершён, результат сохранён в phout/phout${port}_${i}.log"
  done
}

# Запускаем тесты для каждого из портов
run_tests 8081
run_tests 8082

echo "Все тесты завершены."


