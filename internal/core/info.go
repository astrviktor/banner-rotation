package core

// Алгоритм выбора баннера:

// Для каждого банера, который стоит в ротации для заданного слота,
// нужно посчитать величину для заданной группы людей:

// L = xi + sqrt(2 * Ln(n) / ni)

// xi - средний "доход" от i-го баннера для заданной группы людей
// (общее количество переходов для группы / общее количество показов для группы)

// n - общее количество показов всех баннеров из ротации для заданной группы людей
// ni - количество показов i-го баннера из ротации для заданной группы людей

// Исключение: если баннер ни разу не показывался для заданной группы в заданном слоте,
// его надо показать вне очереди для инициализации (при 0 значение формулы стремится к бесконечности)

// Для показа нужно выбрать баннер с максимальным L

// Проверка 1 - нет переходов:

// Есть 1 слот, 2 баннера, группа людей
// 1) Показываем баннер A, 1 показ, 0 переходов (инициализационный показ)
// 2) Показываем баннер B, 1 показ, 0 переходов (инициализационный показ)
// 3)
// Для баннера A: L = 0 + sqrt(2 * ln(2) / 1) =  1,18;
// Для баннера B: L = 0 + sqrt(2 * ln(2) / 1) =  1,18;
//  Выбираем показать A, 2 показа, 0 переходов
// 4)
// Для баннера A: L = 0 + sqrt(2 * ln(3) / 2) =  0,74;
// Для баннера B: L = 0 + sqrt(2 * ln(3) / 1) =  1,48;
//  Выбираем показать B, 2 показа, 0 переходов
// Результат: баннеры меняются по очереди

// Проверка 2 - переход только на B:

// Есть 1 слот, 2 баннера, группа людей
// 1) Показываем баннер A, 1 показ, 0 переходов (инициализационный показ)
// 2) Показываем баннер B, 1 показ, 1 переходов (инициализационный показ)
// 3)
// Для баннера A: L = 0 + sqrt(2 * ln(2) / 1) =  1,18;
// Для баннера B: L = 1 + sqrt(2 * ln(2) / 1) =  2,18;
//  Выбираем показать B, 2 показа, 2 переходов
// 4)
// Для баннера A: L = 0 + sqrt(2 * ln(3) / 1) =  1,48;
// Для баннера B: L = 1 + sqrt(2 * ln(3) / 2) =  1,74;
//  Выбираем показать B, 3 показа, 3 переходов
// 5)
// Для баннера A: L = 0 + sqrt(2 * ln(4) / 1) =  1,67;
// Для баннера B: L = 1 + sqrt(2 * ln(4) / 3) =  1,55;
//  Выбираем показать A, 2 показа, 0 переходов
// 6)
// Для баннера A: L = 0 + sqrt(2 * ln(5) / 2) =  0,89;
// Для баннера B: L = 1 + sqrt(2 * ln(5) / 3) =  1,59;
//  Выбираем показать B, 4 показа, 4 переходов
// 7)
// Для баннера A: L = 0 + sqrt(2 * ln(6) / 2) =  0,94;
// Для баннера B: L = 1 + sqrt(2 * ln(6) / 4) =  1,47;
//  Выбираем показать B, 5 показа, 5 переходов
// 8)
// Для баннера A: L = 0 + sqrt(2 * ln(7) / 2) =  0,98;
// Для баннера B: L = 1 + sqrt(2 * ln(7) / 5) =  1,39;
//  Выбираем показать B, 4 показа, 4 переходов
// Результат: если нет переходов, баннер все равно иногда будет показываться
