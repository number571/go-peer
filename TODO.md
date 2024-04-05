## TODO 

1. Append comments for functions/variables/constants/etc (godoc)
2. Append benchmarks
3. Write linters (names with prefixes 'p', 'f', 's', 'i', 'c', 'g'), (objects with name=p)
4. HLS: write a threat model (in "Hidden Lake" article)
5. Append example of use functions (godoc)
6. Append errors (as traffic/database) for all interface methods
7. Append Batch's to database
8. Custom Unmarshal with fields 
9. golint: "wrapcheck,cyclop,funlen,gocognit"

### Tests

1. Write tests for coverage > 80% (Hidden Lake)
2. Rewrite tests with 'gotests' tool
3. Rewrite http tests with 'httptest' package

### Applications

1. HLFR: Hidden Lake Forum
2. HLN: Hidden Lake Network
3. HLP: Hidden Lake Publisher
4. HLSH: Hidden Lake Shell 

### Articles -> Book

1. Теория строения скрытых систем. Изменено введение в разделе "Анализ сетевой анонимности"
2. Теория строения скрытых систем. Изменено начало в разделе "Анализ сетевой анонимности" подраздела "Стадии анонимности"
3. Теория строения скрытых систем. Дополнена информация в разделе "Анализ сетевой анонимности" подраздела "Стадии анонимности" - "Централизованные сервисы связи (частично) ∈ второй критерий"
4. Абстрактные анонимные сети. Ошибка: "Рисунке 26". Исправить на "Рисунке 2".
5. Теория строения скрытых систем. Текст: Первая^ стадия анонимности `приводит` к необходимости выстраивания большого количества прямых соединений, что `приводит к проблеме` масштабируемости... Изменено `приводит к проблеме` на `порождает проблему`.
6. Децентрализованный протокол обмена ключами. Некорректная формула с расшифрованием: `SPubB = D(PubA, CSPubB)`. Вместо PubA должно быть PrivA. Также необходимо поменять данную формулу на "Рисунке 8". Поменять текст `Пользователь A принимает CSPubB, расшифровывает его своим публичным ключом PubA, получая тем самым SPubB` -> `Пользователь A принимает CSPubB, расшифровывает его своим приватным ключом PrivA, получая тем самым SPubB`.
7. Децентрализованный протокол обмена ключами. Подредактированы пункты 2,3,4 в подразделе "Протокол пятого метода".
