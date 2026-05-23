// Этот файл отвечает только за интерфейс: он не содержит реализации методов.
// Его задача — получать данные из Go-сервера и отрисовавать навигацию.

const taskNav = document.getElementById('taskNav');
const taskTitle = document.getElementById('taskTitle');

const task1Panel = document.getElementById('task1Panel');
const task2Panel = document.getElementById('task2Panel');
const task3Panel = document.getElementById('task3Panel');
const fillPreset3Button = document.getElementById('fillPreset3Button');
const saveTask3InfoButton = document.getElementById('saveTask3InfoButton');
const task3InfoInput = document.getElementById('task3InfoInput');
const task3InfoResult = document.getElementById('task3InfoResult');
const presetSystemView = document.getElementById('presetSystemView');
const fillPresetButton = document.getElementById('fillPresetButton');
const solvePresetButton = document.getElementById('solvePresetButton');
const showGraphButton = document.getElementById('showGraphButton');
const systemSizeInput = document.getElementById('systemSizeInput');
const buildMatrixButton = document.getElementById('buildMatrixButton');
const matrixEditor = document.getElementById('matrixEditor');
const customSystemForm = document.getElementById('customSystemForm');
const epsilonInput = document.getElementById('epsilonInput');
const maxIterationsInput = document.getElementById('maxIterationsInput');
const gaussianResult = document.getElementById('gaussianResult');
const gaussianStepsWrap = document.getElementById('gaussianStepsWrap');
const iterativeResult = document.getElementById('iterativeResult');
const iterationTableWrap = document.getElementById('iterationTableWrap');
const iterationChartWrap = document.getElementById('iterationChartWrap');
const graphSection = document.getElementById('graphSection');
const gaussianCustomResult = document.getElementById('gaussianCustomResult');
const gaussianCustomStepsWrap = document.getElementById('gaussianCustomStepsWrap');
const iterativeCustomResult = document.getElementById('iterativeCustomResult');
const iterationCustomTableWrap = document.getElementById('iterationCustomTableWrap');
const iterationCustomChartWrap = document.getElementById('iterationCustomChartWrap');

const presetSystem2View = document.getElementById('presetSystem2View');
const fillPreset2Button = document.getElementById('fillPreset2Button');
const solvePreset2Button = document.getElementById('solvePreset2Button');
const showGraph2Button = document.getElementById('showGraph2Button');
const sweepPresetResult = document.getElementById('sweepPresetResult');
const sweepPresetIterationTableWrap = document.getElementById('sweepPresetIterationTableWrap');
const sweepPresetIterationChartWrap = document.getElementById('sweepPresetIterationChartWrap');
const seidelPresetResult = document.getElementById('seidelPresetResult');
const seidelPresetIterationTableWrap = document.getElementById('seidelPresetIterationTableWrap');
const seidelPresetIterationChartWrap = document.getElementById('seidelPresetIterationChartWrap');
const graph2PresetSection = document.getElementById('graph2PresetSection');

const system2SizeInput = document.getElementById('system2SizeInput');
const buildMatrix2Button = document.getElementById('buildMatrix2Button');
const matrix2Editor = document.getElementById('matrix2Editor');
const customSystem2Form = document.getElementById('customSystem2Form');
const epsilon2Input = document.getElementById('epsilon2Input');
const maxIterations2Input = document.getElementById('maxIterations2Input');
const sweepCustomResult = document.getElementById('sweepCustomResult');
const sweepCustomIterationTableWrap = document.getElementById('sweepCustomIterationTableWrap');
const sweepCustomIterationChartWrap = document.getElementById('sweepCustomIterationChartWrap');
const seidelCustomResult = document.getElementById('seidelCustomResult');
const seidelCustomIterationTableWrap = document.getElementById('seidelCustomIterationTableWrap');
const seidelCustomIterationChartWrap = document.getElementById('seidelCustomIterationChartWrap');
const graph2CustomSection = document.getElementById('graph2CustomSection');

let tasksCache = [];

// Этот набор данных соответствует варианту 13 из изображения пользователя.
// Он используется как готовый пример и как демонстрация для произвольного ввода.
const presetSystem = {
  matrix: [
    [2.82, 0.43, -0.57],
    [-0.35, 1.12, -0.48],
    [0.48, 0.23, 2.37],
  ],
  vector: [0.48, 0.52, 1.44],
  epsilon: 0.01,
  maxIterations: 100,
  initialGuess: [0, 0, 0],
};

// Исходная система для задания 2 (метод прогонки и метод Зейделя).
const presetSystem2 = {
  matrix: [
    [8, 2, 0, 0],
    [-3, 9, -2, 0],
    [0, 1, 10, 1],
    [0, 0, 1, 6],
  ],
  vector: [15, 5.5, 15, 9.5],
  epsilon: 0.01,
  maxIterations: 100,
  initialGuess: [0, 0, 0, 0],
};

  const presetEquation4 = {
    equation: '3*x - exp(x)',
    a: 0,
    b: 1,
    x0: 1,
    epsilon: 0.0001,
    maxIterations: 100,
  };


// renderTaskButtons отвечает левую панель из данных, пришедших с backend.
function renderTaskButtons(tasks) {
  taskNav.innerHTML = '';

  tasks.forEach((task) => {
    const button = document.createElement('button');
    button.className = 'task-button';
    button.type = 'button';
    button.dataset.id = String(task.id);
    button.innerHTML = `
      <span class="task-button__title">${task.title}</span>
    `;

    button.addEventListener('click', () => {
      // Обновляем hash, чтобы можно было быстро открывать конкретный модуль.
      window.location.hash = `task-${task.id}`;
      activateTask(task.id);
    });

    taskNav.appendChild(button);
  });
}
const task4Panel = document.getElementById('task4Panel');
const fillPreset4Button = document.getElementById('fillPreset4Button');
const solvePreset4Button = document.getElementById('solvePreset4Button');
const bisectionPresetResult = document.getElementById('bisectionPresetResult');
const simpleIterationPresetResult = document.getElementById('simpleIterationPresetResult');
const newtonPresetResult = document.getElementById('newtonPresetResult');
const bisectionPresetIterationTableWrap = document.getElementById('bisectionPresetIterationTableWrap');
const bisectionPresetIterationChartWrap = document.getElementById('bisectionPresetIterationChartWrap');
const simpleIterationPresetIterationTableWrap = document.getElementById('simpleIterationPresetIterationTableWrap');
const simpleIterationPresetIterationChartWrap = document.getElementById('simpleIterationPresetIterationChartWrap');
const newtonPresetIterationTableWrap = document.getElementById('newtonPresetIterationTableWrap');
const newtonPresetIterationChartWrap = document.getElementById('newtonPresetIterationChartWrap');
const customEquation4Form = document.getElementById('customEquation4Form');
const equationInput = document.getElementById('equationInput');
const equation4AInput = document.getElementById('equation4AInput');
const equation4BInput = document.getElementById('equation4BInput');
const equation4X0Input = document.getElementById('equation4X0Input');
const equation4EpsilonInput = document.getElementById('equation4EpsilonInput');
const equation4MaxIterationsInput = document.getElementById('equation4MaxIterationsInput');
const bisectionCustomResult = document.getElementById('bisectionCustomResult');
const simpleIterationCustomResult = document.getElementById('simpleIterationCustomResult');
const newtonCustomResult = document.getElementById('newtonCustomResult');
const bisectionCustomIterationTableWrap = document.getElementById('bisectionCustomIterationTableWrap');
const bisectionCustomIterationChartWrap = document.getElementById('bisectionCustomIterationChartWrap');
const simpleIterationCustomIterationTableWrap = document.getElementById('simpleIterationCustomIterationTableWrap');
const simpleIterationCustomIterationChartWrap = document.getElementById('simpleIterationCustomIterationChartWrap');
const newtonCustomIterationTableWrap = document.getElementById('newtonCustomIterationTableWrap');
const newtonCustomIterationChartWrap = document.getElementById('newtonCustomIterationChartWrap');

// Task 5 elements
const task5Panel = document.getElementById('task5Panel');
const fillPreset5Button = document.getElementById('fillPreset5Button');
const solvePreset5Button = document.getElementById('solvePreset5Button');
const nodesCsvInput = document.getElementById('nodesCsvInput');
const xStarInput = document.getElementById('xStarInput');
const lagrangeResult = document.getElementById('lagrangeResult');
const newtonResult = document.getElementById('newtonResult');
const nodesTableWrap = document.getElementById('nodesTableWrap');
const interpolationChartWrap = document.getElementById('interpolationChartWrap');
const presetSystem5View = document.getElementById('presetSystem5View');

const task6Panel = document.getElementById('task6Panel');
const fillPreset6Button = document.getElementById('fillPreset6Button');
const solvePreset6Button = document.getElementById('solvePreset6Button');
const nodesCsv6Input = document.getElementById('nodesCsv6Input');
const xStar6Input = document.getElementById('xStar6Input');
const boundaryLeftInput = document.getElementById('boundaryLeftInput');
const boundaryRightInput = document.getElementById('boundaryRightInput');
const useTable5Checkbox = document.getElementById('useTable5Checkbox');
const splineValueResult = document.getElementById('splineValueResult');
const splineBoundaryResult = document.getElementById('splineBoundaryResult');
const nodesTable6Wrap = document.getElementById('nodesTable6Wrap');
const splineSegmentsTableWrap = document.getElementById('splineSegmentsTableWrap');
const splineChartWrap = document.getElementById('splineChartWrap');
const presetSystem6View = document.getElementById('presetSystem6View');

const task7Panel = document.getElementById('task7Panel');
const fillPreset7Button = document.getElementById('fillPreset7Button');
const solvePreset7Button = document.getElementById('solvePreset7Button');
const presetSystem7View = document.getElementById('presetSystem7View');
const nodesCsv7Input = document.getElementById('nodesCsv7Input');
const linearResult = document.getElementById('linearResult');
const quadraticResult = document.getElementById('quadraticResult');
const leastSquaresTableWrap = document.getElementById('leastSquaresTableWrap');
const leastSquaresChartWrap = document.getElementById('leastSquaresChartWrap');
// Task 9 elements
const task9Panel = document.getElementById('task9Panel');
const fillPreset9Button = document.getElementById('fillPreset9Button');
const solvePreset9Button = document.getElementById('solvePreset9Button');
const equation9Input = document.getElementById('equation9Input');
const a9Input = document.getElementById('a9Input');
const b9Input = document.getElementById('b9Input');
const eps9Input = document.getElementById('eps9Input');
const integralResult = document.getElementById('integralResult');
const iterationTable9Wrap = document.getElementById('iterationTable9Wrap');
const iterationChart9Wrap = document.getElementById('iterationChart9Wrap');
// Task 10 elements
const task10Panel = document.getElementById('task10Panel');
const fillPreset10Button = document.getElementById('fillPreset10Button');
const solvePreset10Button = document.getElementById('solvePreset10Button');
const equation10Input = document.getElementById('equation10Input');
const a10Input = document.getElementById('a10Input');
const b10Input = document.getElementById('b10Input');
const y0Input = document.getElementById('y0Input');
const h10Input = document.getElementById('h10Input');
const rk4Result = document.getElementById('rk4Result');
const adamsResult = document.getElementById('adamsResult');
const rk4TableWrap = document.getElementById('rk4TableWrap');
const rk4ChartWrap = document.getElementById('rk4ChartWrap');
const adamsTableWrap = document.getElementById('adamsTableWrap');
const adamsChartWrap = document.getElementById('adamsChartWrap');
// Task 11 elements
const task11Panel = document.getElementById('task11Panel');
const fillPreset11Button = document.getElementById('fillPreset11Button');
const solvePreset11Button = document.getElementById('solvePreset11Button');
const p11Input = document.getElementById('p11Input');
const q11Input = document.getElementById('q11Input');
const rhs11Input = document.getElementById('rhs11Input');
const a11Input = document.getElementById('a11Input');
const b11Input = document.getElementById('b11Input');
const h11Input = document.getElementById('h11Input');
const leftAlpha11Input = document.getElementById('leftAlpha11Input');
const leftBeta11Input = document.getElementById('leftBeta11Input');
const leftGamma11Input = document.getElementById('leftGamma11Input');
const rightAlpha11Input = document.getElementById('rightAlpha11Input');
const rightBeta11Input = document.getElementById('rightBeta11Input');
const rightGamma11Input = document.getElementById('rightGamma11Input');
const task11Result = document.getElementById('task11Result');
const task11TableWrap = document.getElementById('task11TableWrap');
const task11ChartWrap = document.getElementById('task11ChartWrap');
// Task 8 elements
const task8Panel = document.getElementById('task8Panel');
const fillPreset8Button = document.getElementById('fillPreset8Button');
const solvePreset8Button = document.getElementById('solvePreset8Button');
const nodesCsv8Input = document.getElementById('nodesCsv8Input');
const x1Input = document.getElementById('x1Input');
const x2Input = document.getElementById('x2Input');
const derivativesResult = document.getElementById('derivativesResult');
const nodesTable8Wrap = document.getElementById('nodesTable8Wrap');
const derivativesChartWrap = document.getElementById('derivativesChartWrap');
// activateTask обновляет центральную область интерфейса.
function activateTask(taskId) {
  const task = tasksCache.find((item) => item.id === taskId);
  if (!task) {
    return;
  }

  document.querySelectorAll('.task-button').forEach((button) => {
    button.classList.toggle('active', Number(button.dataset.id) === taskId);
  });

  taskTitle.textContent = task.title;

  // Показываем только активный модуль.
  task1Panel.classList.toggle('hidden', taskId !== 1);
  task2Panel.classList.toggle('hidden', taskId !== 2);
  if (task3Panel) {
    task3Panel.classList.toggle('hidden', taskId !== 3);
  }
  if (task4Panel) {
    task4Panel.classList.toggle('hidden', taskId !== 4);
  }
  if (typeof task5Panel !== 'undefined' && task5Panel) {
    task5Panel.classList.toggle('hidden', taskId !== 5);
  }
  if (typeof task6Panel !== 'undefined' && task6Panel) {
    task6Panel.classList.toggle('hidden', taskId !== 6);
  }
  if (typeof task7Panel !== 'undefined' && task7Panel) {
    task7Panel.classList.toggle('hidden', taskId !== 7);
  }
  if (typeof task8Panel !== 'undefined' && task8Panel) {
    task8Panel.classList.toggle('hidden', taskId !== 8);
  }
  if (typeof task9Panel !== 'undefined' && task9Panel) {
    task9Panel.classList.toggle('hidden', taskId !== 9);
  }
  if (typeof task10Panel !== 'undefined' && task10Panel) {
    task10Panel.classList.toggle('hidden', taskId !== 10);
  }
  if (typeof task11Panel !== 'undefined' && task11Panel) {
    task11Panel.classList.toggle('hidden', taskId !== 11);
  }
}

// loadTasks получает данные из Go API и поднимает интерфейс.
async function loadTasks() {
  try {
    const response = await fetch('/api/tasks');
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }

    tasksCache = await response.json();
    renderTaskButtons(tasksCache);

    const hash = window.location.hash.replace('#task-', '');
    const taskIdFromHash = Number(hash) || tasksCache[0].id;
    activateTask(taskIdFromHash);
  } catch (error) {
    console.error(error);
    taskTitle.textContent = 'Ошибка загрузки';
  }
}

// buildMatrixEditorTo строит таблицу полей ввода для произвольной СЛАУ.
function buildMatrixEditorTo(editorEl, sizeInputEl, size, system = null, rolePrefix = '') {
  const n = Math.max(2, Math.min(10, Number(size) || 3));
  const matrix = system?.matrix || createEmptyMatrix(n);
  const vector = system?.vector || createEmptyVector(n);

  const table = document.createElement('table');
  table.className = 'matrix-table';

  const thead = document.createElement('thead');
  thead.innerHTML = `<tr><th></th>${Array.from({ length: n }, (_, i) => `<th>a${i + 1}</th>`).join('')}<th>b</th></tr>`;
  table.appendChild(thead);

  const tbody = document.createElement('tbody');
  for (let i = 0; i < n; i++) {
    const row = document.createElement('tr');
    const cells = [`<th>x${i + 1}</th>`];

    for (let j = 0; j < n; j++) {
      const value = Number(matrix[i]?.[j] ?? 0);
      cells.push(`
        <td>
          <input type="number" step="0.01" data-role="${rolePrefix}matrix" data-row="${i}" data-col="${j}" value="${value}">
        </td>
      `);
    }

    const rhs = Number(vector[i] ?? 0);
    cells.push(`
      <td>
        <input type="number" step="0.01" data-role="${rolePrefix}vector" data-row="${i}" value="${rhs}">
      </td>
    `);

    row.innerHTML = cells.join('');
    tbody.appendChild(row);
  }
  table.appendChild(tbody);

  editorEl.innerHTML = '';
  editorEl.appendChild(table);
  sizeInputEl.value = String(n);
  return n;
}

function buildMatrixEditor(size, system = null) {
  return buildMatrixEditorTo(matrixEditor, systemSizeInput, size, system);
}

function buildMatrixEditor2(size, system = null) {
  return buildMatrixEditorTo(matrix2Editor, system2SizeInput, size, system, 'task2-');
}

// createEmptyMatrix создает квадратную матрицу нулей нужного размера.
function createEmptyMatrix(size) {
  return Array.from({ length: size }, () => Array.from({ length: size }, () => 0));
}

// createEmptyVector создает вектор нулей.
function createEmptyVector(size) {
  return Array.from({ length: size }, () => 0);
}

// readSystemFromEditorByRole собирает систему из указанной формы.
function readSystemFromEditorByRole(editorEl, sizeInputEl, epsilonEl, maxIterationsEl, rolePrefix = '') {
  const matrixInputs = Array.from(editorEl.querySelectorAll(`input[data-role="${rolePrefix}matrix"]`));
  const vectorInputs = Array.from(editorEl.querySelectorAll(`input[data-role="${rolePrefix}vector"]`));
  const size = Number(sizeInputEl.value) || 3;
  const matrix = createEmptyMatrix(size);
  const vector = createEmptyVector(size);

  matrixInputs.forEach((input) => {
    const row = Number(input.dataset.row);
    const col = Number(input.dataset.col);
    matrix[row][col] = Number(input.value);
  });

  vectorInputs.forEach((input) => {
    const row = Number(input.dataset.row);
    vector[row] = Number(input.value);
  });

  return {
    matrix,
    vector,
    epsilon: Number(epsilonEl.value) || 0.01,
    maxIterations: Number(maxIterationsEl.value) || 100,
    initialGuess: createEmptyVector(size),
  };
}

function readSystemFromEditor() {
  return readSystemFromEditorByRole(matrixEditor, systemSizeInput, epsilonInput, maxIterationsInput);
}

function readSystem2FromEditor() {
  return readSystemFromEditorByRole(matrix2Editor, system2SizeInput, epsilon2Input, maxIterations2Input, 'task2-');
}

// formatVector превращает вектор в читаемую строку.
function formatVector(values) {
  if (!Array.isArray(values) || values.length === 0) {
    return 'Нет данных';
  }

  return values.map((value, index) => `x${index + 1} = ${Number(value).toFixed(6)}`).join('<br>');
}

// formatNumeric выводит число либо в fixed с 6 знаками, либо в экспоненциальной
// форме, если значение по абсолютной величине очень мало (чтобы не показывать "0.000000").
function formatNumeric(v) {
  if (typeof v !== 'number' || Number.isNaN(v)) return '0.000000';
  const abs = Math.abs(v);
  if (abs === 0) return '0.000000';
  if (abs < 1e-6) return v.toExponential(3);
  return v.toFixed(6);
}

// renderResultList показывает список сообщений в карточке результата.
function renderResultList(container, values, title) {
  if (!Array.isArray(values) || values.length === 0) {
    container.innerHTML = '<span class="chart-empty">Результат отсутствует.</span>';
    return;
  }

  container.innerHTML = `
    <strong>${title}</strong>
    <ul class="result-list">
      ${values.map((item, index) => `<li>${index + 1}. ${item}</li>`).join('')}
    </ul>
  `;
}

function buildMethodSummaryLines(result, options = {}) {
  const lines = [];
  const errorText = options.errorText || 'Ошибка расчета.';
  const notConvergedText = options.notConvergedText || 'Метод не сошелся.';
  const convergedText = options.convergedText || 'Метод сошелся.';
  const residualLabel = options.residualLabel || 'Последняя невязка |f(x)|';
  const solutionLabel = options.solutionLabel || 'Решение';

  if (result?.error) {
    lines.push(result.error);
  }

  lines.push(result?.converged ? convergedText : notConvergedText);

  const delta = result?.finalDelta;
  const residual = result?.finalResidual;
  if (typeof delta === 'number') {
    lines.push(`Последняя δ = ${formatNumeric(delta)}`);
  } else {
    lines.push('Последняя δ = 0.000000');
  }
  if (typeof residual === 'number') {
    lines.push(`${residualLabel} = ${formatNumeric(residual)}`);
  } else {
    lines.push(`${residualLabel} = 0.000000`);
  }

  const solution = Array.isArray(result?.solution) ? result.solution : null;
  if (solution && solution.length > 0) {
    lines.push(`${solutionLabel}: ${formatVector(solution).replaceAll('<br>', '; ')}`);
  } else {
    lines.push(`${solutionLabel}: нет данных`);
  }

  return lines;
}

// renderTask1ResponseIn заполняет блоки результата для Task1 в указанные контейнеры.
function renderTask1ResponseIn(response, target) {
  const gaussian = response.gaussian;
  const iterative = response.iterative;
  const iterations = iterative.iterations || [];

  if (gaussian.error) {
    target.gaussianResult.innerHTML = `<span class="chart-empty">${gaussian.error}</span>`;
    target.stepsWrap.innerHTML = '<span class="chart-empty">Шаги преобразования недоступны.</span>';
  } else {
    const gaussianLines = buildMethodSummaryLines(gaussian, {
      convergedText: 'Метод сошелся.',
      notConvergedText: 'Метод не сошелся.',
      residualLabel: 'Последняя невязка |f(x)|',
      solutionLabel: 'Решение',
    });
    renderResultList(target.gaussianResult, gaussianLines, 'Метод Гаусса');
    renderGaussianStepsIn(target.stepsWrap, gaussian.steps || []);
  }

  const iterativeLines = buildMethodSummaryLines(iterative, {
    convergedText: 'Метод сошелся.',
    notConvergedText: 'Метод не сошелся.',
    residualLabel: 'Последняя невязка |f(x)|',
    solutionLabel: 'Решение',
  });

  renderResultList(target.iterativeResult, iterativeLines, 'Метод простой итерации');
  renderIterationTableIn(target.iterationTableWrap, iterations);
  renderIterationChartIn(target.iterationChartWrap, iterations);
}

// renderSolveResponse заполняет все блоки результата для исходной системы.
function renderSolveResponse(response) {
  renderTask1ResponseIn(response, {
    gaussianResult,
    stepsWrap: gaussianStepsWrap,
    iterativeResult,
    iterationTableWrap,
    iterationChartWrap,
  });
}

// renderGaussianStepsIn выводит пошаговый ход преобразования матрицы Гаусса в контейнер.
function renderGaussianStepsIn(container, steps) {
  if (!Array.isArray(steps) || steps.length === 0) {
    container.innerHTML = '<span class="chart-empty">Шаги преобразования не найдены.</span>';
    return;
  }

  container.innerHTML = steps.map((step, index) => {
    const matrixHtml = renderMatrixTable(step.matrix);
    return `
      <article class="gaussian-step">
        <div class="gaussian-step__head">
          <span class="gaussian-step__index">Шаг ${index + 1}</span>
          <strong>${step.title}</strong>
        </div>
        <div class="gaussian-step__body">
          ${matrixHtml}
        </div>
      </article>
    `;
  }).join('');
}

// renderGaussianSteps выводит пошаговый ход преобразования матрицы Гаусса.
function renderGaussianSteps(steps) {
  renderGaussianStepsIn(gaussianStepsWrap, steps);
}

// renderMatrixTable превращает расширенную матрицу в компактную HTML-таблицу.
function renderMatrixTable(matrix) {
  if (!Array.isArray(matrix) || matrix.length === 0) {
    return '<span class="chart-empty">Матрица недоступна.</span>';
  }

  const columnCount = matrix[0].length;
  const variableCount = Math.max(0, columnCount - 1);
  const header = [
    ...Array.from({ length: variableCount }, (_, i) => `<th>a${i + 1}</th>`),
    '<th>b</th>',
  ].join('');

  const rows = matrix.map((row) => {
    return `<tr>${row.map((value) => `<td>${Number(value).toFixed(6)}</td>`).join('')}</tr>`;
  }).join('');

  return `
    <table class="gaussian-matrix">
      <thead>
        <tr>${header}</tr>
      </thead>
      <tbody>${rows}</tbody>
    </table>
  `;
}

// renderIterationTableIn строит таблицу по всем итерациям в нужном контейнере.
function renderIterationTableIn(container, iterations) {
  if (!Array.isArray(iterations) || iterations.length === 0) {
    container.innerHTML = '<span class="chart-empty">Итерации пока отсутствуют.</span>';
    return;
  }

  const firstRow = iterations[0] || {};
  const hasVectorValues = Array.isArray(firstRow.values);
  const variableCount = hasVectorValues ? firstRow.values.length : (typeof firstRow.x === 'number' ? 1 : 0);
  const headerCells = [
    '<th>#</th>',
    ...Array.from({ length: variableCount }, (_, i) => `<th>x${i + 1}</th>`),
    '<th>Δ</th>',
    '<th>Невязка</th>',
  ].join('');

  const rows = iterations.map((row) => {
    const values = hasVectorValues
      ? row.values.map((value) => `<td>${Number(value).toFixed(6)}</td>`).join('')
      : `<td>${Number(row.x ?? 0).toFixed(6)}</td>`;
    return `
      <tr>
        <td>${row.iteration}</td>
        ${values}
        <td>${formatNumeric(Number(row.delta) || 0)}</td>
        <td>${formatNumeric(Number(row.residual) || 0)}</td>
      </tr>
    `;
  }).join('');

  container.innerHTML = `
    <table class="data-table">
      <thead>
        <tr>${headerCells}</tr>
      </thead>
      <tbody>${rows}</tbody>
    </table>
  `;
}

function renderIterationTable(iterations) {
  renderIterationTableIn(iterationTableWrap, iterations);
}

// renderIterationChartIn строит SVG-график по невязке на каждой итерации.
function renderIterationChartIn(container, iterations) {
  if (!Array.isArray(iterations) || iterations.length === 0) {
    container.innerHTML = '<span class="chart-empty">График появится после расчета.</span>';
    return;
  }

  const values = iterations.map((item) => Number(item.residual));
  const width = 900;
  const height = 320;
  const padding = 48;

  // Если есть только одна точка, сделаем запас сверху/снизу около 10% от величины,
  // чтобы точка не была прижата к оси и была визуально заметна.
  const rawMax = Math.max(...values);
  const rawMin = Math.min(...values);
  let minValue = 0;
  let maxValue = Math.max(rawMax, 1e-12);
  if (values.length === 1) {
    const v = values[0] || 0;
    const abs = Math.abs(v);
    const delta = Math.max(abs * 0.1, 1e-12);
    minValue = v - delta;
    maxValue = v + delta;
  } else {
    // при нескольких точках ставим нижний предел в 0 для визуальной совместимости
    minValue = 0;
    maxValue = Math.max(rawMax, 1e-6);
  }

  const stepX = values.length === 1 ? 0 : (width - padding * 2) / (values.length - 1);
  const usableHeight = height - padding * 2;

  const points = values.map((value, index) => {
    const x = padding + index * stepX;
    const normalized = maxValue === minValue ? 0.5 : (value - minValue) / (maxValue - minValue);
    const y = height - padding - normalized * usableHeight;
    return { x, y, value };
  });

  const polyline = points.map((point) => `${point.x},${point.y}`).join(' ');

  const ticks = Array.from({ length: 5 }, (_, index) => {
    const ratio = index / 4;
    const value = minValue + (1 - ratio) * (maxValue - minValue);
    const y = padding + ratio * usableHeight;
    return `
      <g>
        <line x1="${padding}" y1="${y}" x2="${width - padding}" y2="${y}" stroke="rgba(255,255,255,0.08)" />
        <text x="12" y="${y + 4}" fill="#98a8c7" font-size="12">${formatNumeric(value)}</text>
      </g>
    `;
  }).join('');

  const xLabels = points.map((point, index) => `
    <text x="${point.x}" y="${height - 18}" text-anchor="middle" fill="#98a8c7" font-size="12">${index + 1}</text>
  `).join('');

  container.innerHTML = `
    <svg viewBox="0 0 ${width} ${height}" class="chart-svg" role="img" aria-label="График невязки по итерациям">
      ${ticks}
      <line x1="${padding}" y1="${padding}" x2="${padding}" y2="${height - padding}" stroke="rgba(255,255,255,0.18)" />
      <line x1="${padding}" y1="${height - padding}" x2="${width - padding}" y2="${height - padding}" stroke="rgba(255,255,255,0.18)" />
      <polyline fill="none" stroke="#67e8f9" stroke-width="3" points="${polyline}" />
      ${points.map((point) => `<circle cx="${point.x}" cy="${point.y}" r="4" fill="#8b5cf6"><title>Итерация ${formatNumeric(point.value)}</title></circle>`).join('')}
      ${xLabels}
    </svg>
  `;
}

function renderIterationChart(iterations) {
  renderIterationChartIn(iterationChartWrap, iterations);
}

// fillPresetSystem заполняет форму данными варианта 13.
function fillPresetSystem() {
  buildMatrixEditor(presetSystem.matrix.length, presetSystem);
  epsilonInput.value = String(presetSystem.epsilon);
  maxIterationsInput.value = String(presetSystem.maxIterations);
  presetSystemView.textContent = presetSystem.matrix
    .map((row, index) => `${row.map((value, j) => `${value >= 0 ? ' ' : ''}${value.toFixed(2)}x${j + 1}`).join(' + ')} = ${presetSystem.vector[index].toFixed(2)}`)
    .join('\n');
}

// solvePreset отправляет на сервер готовый вариант из задания.
async function solvePreset() {
  await solveSystem(presetSystem);
}

// solveCustom1 отправляет произвольную СЛАУ на сервер и рендерит результаты в custom контейнеры.
async function solveCustom1() {
  const system = readSystemFromEditor();
  await solveTask1System(system, {
    gaussianResult: gaussianCustomResult,
    stepsWrap: gaussianCustomStepsWrap,
    iterativeResult: iterativeCustomResult,
    iterationTableWrap: iterationCustomTableWrap,
    iterationChartWrap: iterationCustomChartWrap,
  });
}

// showGraph плавно прокручивает страницу к блоку графика.
function showGraph() {
  graphSection?.scrollIntoView({ behavior: 'smooth', block: 'start' });
}

// solveSystem отправляет СЛАУ на backend и получает оба решения.
async function solveTask1System(payload, target) {
  try {
    const response = await fetch('/api/task1/solve', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }

    const data = await response.json();
    renderTask1ResponseIn(data, target);
  } catch (error) {
    console.error(error);
    target.gaussianResult.innerHTML = '<span class="chart-empty">Ошибка расчета.</span>';
    target.stepsWrap.innerHTML = '<span class="chart-empty">Шаги недоступны.</span>';
    target.iterativeResult.innerHTML = `<span class="chart-empty">${String(error)}</span>`;
    target.iterationTableWrap.innerHTML = '<span class="chart-empty">Таблица недоступна.</span>';
    target.iterationChartWrap.innerHTML = '<span class="chart-empty">График недоступен.</span>';
  }
}

async function solveSystem(system = null) {
  const payload = system || readSystemFromEditor();

  await solveTask1System(payload, {
    gaussianResult,
    stepsWrap: gaussianStepsWrap,
    iterativeResult,
    iterationTableWrap,
    iterationChartWrap,
  });
}

function renderTask2Response(response, target) {
  const sweep = response.sweep || {};
  const seidel = response.seidel || {};
  const sweepIterations = sweep.iterations || [];
  const iterations = seidel.iterations || [];

  const sweepLines = buildMethodSummaryLines(sweep, {
    convergedText: 'Метод сошелся.',
    notConvergedText: 'Метод не сошелся.',
    residualLabel: 'Последняя невязка |f(x)|',
    solutionLabel: 'Решение',
  });
  renderResultList(target.sweepResult, sweepLines, 'Метод прогонки');
  renderIterationTableIn(target.sweepIterationTableWrap, sweepIterations);
  renderIterationChartIn(target.sweepIterationChartWrap, sweepIterations);

  const seidelLines = buildMethodSummaryLines(seidel, {
    convergedText: 'Метод сошелся.',
    notConvergedText: 'Метод не сошелся.',
    residualLabel: 'Последняя невязка |f(x)|',
    solutionLabel: 'Решение',
  });
  renderResultList(target.seidelResult, seidelLines, 'Метод Зейделя');
  renderIterationTableIn(target.iterationTableWrap, iterations);
  renderIterationChartIn(target.iterationChartWrap, iterations);
}

async function solveTask2System(payload, target) {
  try {
    const response = await fetch('/api/task2/solve', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }

    const data = await response.json();
    renderTask2Response(data, target);
  } catch (error) {
    console.error(error);
    target.sweepResult.innerHTML = '<span class="chart-empty">Ошибка расчета.</span>';
    target.seidelResult.innerHTML = `<span class="chart-empty">${String(error)}</span>`;
    target.sweepIterationTableWrap.innerHTML = '<span class="chart-empty">Таблица недоступна.</span>';
    target.sweepIterationChartWrap.innerHTML = '<span class="chart-empty">График недоступен.</span>';
    target.iterationTableWrap.innerHTML = '<span class="chart-empty">Таблица недоступна.</span>';
    target.iterationChartWrap.innerHTML = '<span class="chart-empty">График недоступен.</span>';
  }
}

function fillPresetSystem2() {
  buildMatrixEditor2(presetSystem2.matrix.length, presetSystem2);
  epsilon2Input.value = String(presetSystem2.epsilon);
  maxIterations2Input.value = String(presetSystem2.maxIterations);
  presetSystem2View.textContent = [
    '8.00x1 + 2.00x2 = 15.00',
    '-3.00x1 + 9.00x2 - 2.00x3 = 5.50',
    '1.00x2 + 10.00x3 + 1.00x4 = 15.00',
    '1.00x3 + 6.00x4 = 9.50',
  ].join('\n');
}

async function solvePreset2() {
  await solveTask2System(presetSystem2, {
    sweepResult: sweepPresetResult,
    sweepIterationTableWrap: sweepPresetIterationTableWrap,
    sweepIterationChartWrap: sweepPresetIterationChartWrap,
    seidelResult: seidelPresetResult,
    iterationTableWrap: seidelPresetIterationTableWrap,
    iterationChartWrap: seidelPresetIterationChartWrap,
  });
}

async function solveCustom2() {
  await solveTask2System(readSystem2FromEditor(), {
    sweepResult: sweepCustomResult,
    sweepIterationTableWrap: sweepCustomIterationTableWrap,
    sweepIterationChartWrap: sweepCustomIterationChartWrap,
    seidelResult: seidelCustomResult,
    iterationTableWrap: seidelCustomIterationTableWrap,
    iterationChartWrap: seidelCustomIterationChartWrap,
  });
}

function showGraph2() {
  graph2PresetSection?.scrollIntoView({ behavior: 'smooth', block: 'start' });
}

// Если пользователь меняет hash вручную, интерфейс тоже должен обновиться.
window.addEventListener('hashchange', () => {
  const hash = window.location.hash.replace('#task-', '');
  const taskId = Number(hash);
  if (!Number.isNaN(taskId)) {
    activateTask(taskId);
  }
});

// Task 4: Решение нелинейных уравнений
function fillPresetEquation4() {
  equationInput.value = presetEquation4.equation;
  equation4AInput.value = presetEquation4.a;
  equation4BInput.value = presetEquation4.b;
  equation4X0Input.value = presetEquation4.x0;
  equation4EpsilonInput.value = presetEquation4.epsilon;
  equation4MaxIterationsInput.value = presetEquation4.maxIterations;
  presetSystem4View.textContent = `${presetEquation4.equation} = 0 на [${presetEquation4.a}, ${presetEquation4.b}]`;
}

function fillPreset3Info() {
  if (task3InfoInput) {
    task3InfoInput.value = 'Заметка по заданию 3:\n- Введите свои исходные данные\n- Сохраните, чтобы видеть текст ниже';
  }
  if (task3InfoResult) {
    task3InfoResult.innerHTML = '<span class="chart-empty">Пока ничего не сохранено.</span>';
  }
}

function saveTask3Info() {
  if (!task3InfoResult) {
    return;
  }
  const text = (task3InfoInput && task3InfoInput.value ? task3InfoInput.value.trim() : '');
  if (!text) {
    task3InfoResult.innerHTML = '<span class="chart-empty">Введите текст перед сохранением.</span>';
    return;
  }
  task3InfoResult.textContent = text;
}

function fillPreset5() {
  // оставляем значение из textarea по умолчанию, просто обновим превью
  if (presetSystem5View) {
    presetSystem5View.textContent = 'Введите узлы в CSV (xi, yi) по строкам. Пример заполнен.';
  }
}

function fillPreset6() {
  if (presetSystem6View) {
    presetSystem6View.textContent = 'Натуральный кубический сплайн по заданным узлам.';
  }
  if (useTable5Checkbox) {
    useTable5Checkbox.checked = true;
  }
  if (boundaryLeftInput) {
    boundaryLeftInput.value = '';
  }
  if (boundaryRightInput) {
    boundaryRightInput.value = '';
  }
}

function evaluateSplineAt(segments, x) {
  if (!Array.isArray(segments) || segments.length === 0) {
    return 0;
  }

  let segment = segments[segments.length - 1];
  for (const candidate of segments) {
    if (x >= candidate.x0 && x <= candidate.x1) {
      segment = candidate;
      break;
    }
  }

  if (x < segments[0].x0) {
    segment = segments[0];
  }

  const dx = x - segment.x0;
  return segment.a + segment.b * dx + segment.c * dx * dx + segment.d * dx * dx * dx;
}

function renderSplineSegmentsTable(container, segments) {
  if (!Array.isArray(segments) || segments.length === 0) {
    container.innerHTML = '<span class="chart-empty">Сегменты сплайна отсутствуют.</span>';
    return;
  }

  const rows = segments.map((segment) => `
    <tr>
      <td>${segment.index}</td>
      <td>[${Number(segment.x0).toFixed(6)}, ${Number(segment.x1).toFixed(6)}]</td>
      <td>S${segment.index}(x) = ${Number(segment.a).toFixed(6)} + ${Number(segment.b).toFixed(6)}(x-x${segment.index}) + ${Number(segment.c).toFixed(6)}(x-x${segment.index})² + ${Number(segment.d).toFixed(6)}(x-x${segment.index})³</td>
    </tr>
  `).join('');

  container.innerHTML = `
    <table class="data-table">
      <thead>
        <tr>
          <th>#</th>
          <th>Отрезок</th>
          <th>Кубический сплайн</th>
        </tr>
      </thead>
      <tbody>${rows}</tbody>
    </table>
  `;
}

function renderSplineChart(container, nodes, segments, evaluations = []) {
  if (!Array.isArray(nodes) || nodes.length === 0 || !Array.isArray(segments) || segments.length === 0) {
    container.innerHTML = '<span class="chart-empty">График недоступен.</span>';
    return;
  }

  const xs = nodes.map((node) => node.x);
  const ys = nodes.map((node) => node.y);
  const minX = Math.min(...xs);
  const maxX = Math.max(...xs);
  const margin = (maxX - minX) * 0.08 || 1;
  const samples = 220;
  const curve = [];

  for (let i = 0; i <= samples; i++) {
    const x = minX - margin + (i / samples) * (maxX - minX + 2 * margin);
    curve.push({ x, y: evaluateSplineAt(segments, x) });
  }

  const allY = curve.map((point) => point.y).concat(ys);
  const minY = Math.min(...allY);
  const maxY = Math.max(...allY);

  const width = 900;
  const height = 320;
  const padding = 48;
  const usableW = width - padding * 2;
  const usableH = height - padding * 2;
  const rangeX = maxX - minX + 2 * margin || 1;
  const rangeY = maxY - minY || 1;

  const toX = (x) => padding + ((x - (minX - margin)) / rangeX) * usableW;
  const toY = (y) => height - padding - ((y - minY) / rangeY) * usableH;

  const polyline = curve.map((point) => `${toX(point.x)},${toY(point.y)}`).join(' ');
  const nodeCircles = nodes.map((node) => `<circle cx="${toX(node.x)}" cy="${toY(node.y)}" r="4" fill="#f97316"><title>(${node.x.toFixed(6)}, ${node.y.toFixed(6)})</title></circle>`).join('');
  const evalCircles = Array.isArray(evaluations) ? evaluations.map((e) => `<circle cx="${toX(e.x0)}" cy="${toY(e.value)}" r="3" fill="#f43f5e"><title>(${e.x0.toFixed(6)}, ${e.value.toFixed(6)})</title></circle>`).join('') : '';

  // axis ticks
  const xStart = minX - margin;
  const xRange = rangeX;
  const xTicks = Array.from({ length: 5 }, (_, i) => {
    const rx = i / 4;
    const xv = xStart + rx * xRange;
    return { x: toX(xv), v: xv };
  });
  const yTicks = Array.from({ length: 5 }, (_, i) => {
    const ry = i / 4;
    const yv = minY + (1 - ry) * (maxY - minY);
    return { y: toY(yv), v: yv };
  });

  const xTickLines = xTicks.map((t) => `<text x="${t.x}" y="${height - 14}" text-anchor="middle" fill="#98a8c7" font-size="12">${t.v.toFixed(4)}</text>`).join('');
  const yTickLines = yTicks.map((t) => `<g><line x1="${padding}" y1="${t.y}" x2="${width - padding}" y2="${t.y}" stroke="rgba(255,255,255,0.06)" /><text x="12" y="${t.y + 4}" fill="#98a8c7" font-size="12">${t.v.toFixed(4)}</text></g>`).join('');

  container.innerHTML = `
    <svg viewBox="0 0 ${width} ${height}" class="chart-svg" role="img" aria-label="График кубического сплайна">
      ${yTickLines}
      <line x1="${padding}" y1="${padding}" x2="${padding}" y2="${height - padding}" stroke="rgba(255,255,255,0.18)" />
      <line x1="${padding}" y1="${height - padding}" x2="${width - padding}" y2="${height - padding}" stroke="rgba(255,255,255,0.18)" />
      ${xTickLines}
      <polyline fill="none" stroke="#38bdf8" stroke-width="3" points="${polyline}" />
      ${nodeCircles}
      ${evalCircles}
    </svg>
  `;
}

async function solveTask6() {
  const csv = nodesCsv6Input.value || '';
  const xStar = parseFloat(xStar6Input.value) || 0;
  try {
    const payload = { csv, x_star: xStar };
    const leftStr = (boundaryLeftInput && boundaryLeftInput.value) ? boundaryLeftInput.value.trim() : '';
    const rightStr = (boundaryRightInput && boundaryRightInput.value) ? boundaryRightInput.value.trim() : '';
    if (leftStr !== '') payload.boundary_left_second = parseFloat(leftStr);
    if (rightStr !== '') payload.boundary_right_second = parseFloat(rightStr);
    if (useTable5Checkbox && useTable5Checkbox.checked) {
      // parse nodesCsv6Input to get x values
      const lines = (nodesCsv6Input.value || '').split('\n').map((s) => s.trim()).filter(Boolean);
      const pts = [];
      for (const ln of lines) {
        const parts = ln.split(',');
        if (parts.length >= 1) {
          const x = parseFloat(parts[0]);
          if (!Number.isNaN(x)) pts.push(x);
        }
      }
      if (pts.length > 0) payload.evaluate_points = pts;
    }

    const response = await fetch('/api/task6/solve', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }

    const data = await response.json();
    renderTask6Response(data);
  } catch (error) {
    console.error(error);
    splineValueResult.innerHTML = `<span class="chart-empty">Ошибка: ${String(error)}</span>`;
    splineBoundaryResult.innerHTML = '<span class="chart-empty">Ошибка расчета.</span>';
    splineSegmentsTableWrap.innerHTML = '<span class="chart-empty">Таблица недоступна.</span>';
    splineChartWrap.innerHTML = '<span class="chart-empty">График недоступен.</span>';
  }
}

function renderTask6Response(response) {
  const nodes = response.nodes || [];
  const segments = response.segments || [];
  const evals = response.evaluations || [];
  let valueHtml = '';
  if (Array.isArray(evals) && evals.length > 0) {
    valueHtml += '<table class="data-table"><thead><tr><th>#</th><th>x</th><th>S(x)</th></tr></thead><tbody>';
    for (let i = 0; i < evals.length; i++) {
      const e = evals[i];
      const x = Number(e.x0 ?? 0);
      const v = Number(e.value ?? 0);
      valueHtml += `<tr><td>${i+1}</td><td>${x.toFixed(6)}</td><td>${v.toFixed(6)}</td></tr>`;
    }
    valueHtml += '</tbody></table>';
  } else {
    valueHtml = `<div>Значение сплайна = ${Number(response.spline_value).toFixed(6)}</div>`;
  }

  splineValueResult.innerHTML = valueHtml;
  splineBoundaryResult.innerHTML = `<div>S''(x1) = ${Number(response.boundary_left_second).toFixed(6)}<br>S''(xn) = ${Number(response.boundary_right_second).toFixed(6)}</div>`;
  renderNodesTable(nodesTable6Wrap, nodes);
  renderSplineSegmentsTable(splineSegmentsTableWrap, segments);
  renderSplineChart(splineChartWrap, nodes, segments, evals);
}

const task7PresetPoints = [
  { x: 0.0, y: 1.0 },
  { x: 0.12, y: 1.2 },
  { x: 0.19, y: 1.6 },
  { x: 0.35, y: 2.6 },
  { x: 0.4, y: 1.8 },
  { x: 0.45, y: 2.7 },
  { x: 0.62, y: 3.5 },
  { x: 0.71, y: 4.4 },
  { x: 0.84, y: 4.5 },
  { x: 0.91, y: 5.2 },
  { x: 1.0, y: 6.3 },
];

function formatTask7Table(points) {
  return points.map((point) => `${point.x.toFixed(2)}\t${point.y.toFixed(2)}`).join('\n');
}

function fillPreset7() {
  if (presetSystem7View) {
    presetSystem7View.textContent = formatTask7Table(task7PresetPoints);
  }
  if (nodesCsv7Input) {
    nodesCsv7Input.value = formatTask7Table(task7PresetPoints);
  }
}

function parseCsvPoints(text) {
  const lines = String(text || '').split(/\r?\n/).map(l => l.trim()).filter(l => l.length > 0);
  const points = [];
  for (const line of lines) {
    // split by comma, semicolon or whitespace (tabs allowed)
    const parts = line.split(/[,;\s]+/).map(p => p.trim()).filter(p => p.length > 0);
    if (parts.length < 2) continue;
    const x = parseFloat(parts[0]);
    const y = parseFloat(parts[1]);
    if (Number.isNaN(x) || Number.isNaN(y)) {
      throw new Error(`Неверный формат строки: "${line}"`);
    }
    points.push({ x, y });
  }
  if (points.length === 0) throw new Error('Нет корректных точек для расчета.');
  return points;
}

function renderLeastSquaresResults(response) {
  const linear = response.linear || {};
  const quadratic = response.quadratic || {};
  const rows = response.rows || [];

  linearResult.innerHTML = `
    <div>y = ${Number(linear.a).toFixed(6)} + ${Number(linear.b).toFixed(6)}x</div>
    <div>RMSE = ${Number(linear.rmse).toFixed(6)}</div>
    <div>SSE = ${Number(linear.sse).toFixed(6)}</div>
  `;

  quadraticResult.innerHTML = `
    <div>y = ${Number(quadratic.a).toFixed(6)} + ${Number(quadratic.b).toFixed(6)}x + ${Number(quadratic.c).toFixed(6)}x²</div>
    <div>RMSE = ${Number(quadratic.rmse).toFixed(6)}</div>
    <div>SSE = ${Number(quadratic.sse).toFixed(6)}</div>
  `;

  const tableRows = rows.map((row) => `
    <tr>
      <td>${row.index}</td>
      <td>${Number(row.x).toFixed(2)}</td>
      <td>${Number(row.y).toFixed(2)}</td>
      <td>${Number(row.linear).toFixed(6)}</td>
      <td>${Number(row.quadratic).toFixed(6)}</td>
      <td>${Number(row.linear_error).toFixed(6)}</td>
      <td>${Number(row.quadratic_error).toFixed(6)}</td>
    </tr>
  `).join('');

  leastSquaresTableWrap.innerHTML = `
    <table class="data-table">
      <thead>
        <tr>
          <th>#</th>
          <th>x</th>
          <th>y</th>
          <th>Линейная</th>
          <th>Квадратичная</th>
          <th>Ошибка линейной</th>
          <th>Ошибка квадратичной</th>
        </tr>
      </thead>
      <tbody>${tableRows}</tbody>
    </table>
  `;

  renderLeastSquaresChart(response.points || task7PresetPoints, linear, quadratic);
}

function renderLeastSquaresChart(points, linear, quadratic) {
  if (!Array.isArray(points) || points.length === 0) {
    leastSquaresChartWrap.innerHTML = '<span class="chart-empty">График недоступен.</span>';
    return;
  }

  const xs = points.map((point) => point.x);
  const ys = points.map((point) => point.y);
  const minX = Math.min(...xs);
  const maxX = Math.max(...xs);
  const margin = (maxX - minX) * 0.1 || 1;
  const samples = 220;
  const linearCurve = [];
  const quadraticCurve = [];

  for (let i = 0; i <= samples; i++) {
    const x = minX - margin + (i / samples) * (maxX - minX + 2 * margin);
    linearCurve.push({ x, y: linear.a + linear.b * x });
    quadraticCurve.push({ x, y: quadratic.a + quadratic.b * x + quadratic.c * x * x });
  }

  const allY = ys.concat(linearCurve.map((point) => point.y), quadraticCurve.map((point) => point.y));
  const minY = Math.min(...allY);
  const maxY = Math.max(...allY);

  const width = 900;
  const height = 320;
  const padding = 48;
  const usableW = width - padding * 2;
  const usableH = height - padding * 2;
  const rangeX = maxX - minX + 2 * margin || 1;
  const rangeY = maxY - minY || 1;

  const toX = (x) => padding + ((x - (minX - margin)) / rangeX) * usableW;
  const toY = (y) => height - padding - ((y - minY) / rangeY) * usableH;

  const pointCircles = points.map((point) => `<circle cx="${toX(point.x)}" cy="${toY(point.y)}" r="4" fill="#f97316"><title>(${point.x.toFixed(2)}, ${point.y.toFixed(2)})</title></circle>`).join('');
  const linearPolyline = linearCurve.map((point) => `${toX(point.x)},${toY(point.y)}`).join(' ');
  const quadraticPolyline = quadraticCurve.map((point) => `${toX(point.x)},${toY(point.y)}`).join(' ');

  // axis ticks
  const xStart = minX - margin;
  const xRange = rangeX;
  const xTicks = Array.from({ length: 5 }, (_, i) => {
    const rx = i / 4;
    const xv = xStart + rx * xRange;
    return { x: toX(xv), v: xv };
  });
  const yTicks = Array.from({ length: 5 }, (_, i) => {
    const ry = i / 4;
    const yv = minY + (1 - ry) * (maxY - minY);
    return { y: toY(yv), v: yv };
  });

  const xTickLines = xTicks.map((t) => `<text x="${t.x}" y="${height - 14}" text-anchor="middle" fill="#98a8c7" font-size="12">${t.v.toFixed(4)}</text>`).join('');
  const yTickLines = yTicks.map((t) => `<g><line x1="${padding}" y1="${t.y}" x2="${width - padding}" y2="${t.y}" stroke="rgba(255,255,255,0.06)" /><text x="12" y="${t.y + 4}" fill="#98a8c7" font-size="12">${t.v.toFixed(4)}</text></g>`).join('');

  leastSquaresChartWrap.innerHTML = `
    <svg viewBox="0 0 ${width} ${height}" class="chart-svg" role="img" aria-label="График аппроксимации МНК">
      ${yTickLines}
      <line x1="${padding}" y1="${padding}" x2="${padding}" y2="${height - padding}" stroke="rgba(255,255,255,0.18)" />
      <line x1="${padding}" y1="${height - padding}" x2="${width - padding}" y2="${height - padding}" stroke="rgba(255,255,255,0.18)" />
      ${xTickLines}
      <polyline fill="none" stroke="#67e8f9" stroke-width="3" points="${linearPolyline}" />
      <polyline fill="none" stroke="#a78bfa" stroke-width="3" points="${quadraticPolyline}" />
      ${pointCircles}
    </svg>
  `;
}

function renderLineChart(container, seriesArray) {
  if (!Array.isArray(seriesArray) || seriesArray.length === 0) {
    container.innerHTML = '<span class="chart-empty">График недоступен.</span>';
    return;
  }

  const allPoints = seriesArray.flatMap(s => s.points || []);
  if (allPoints.length === 0) {
    container.innerHTML = '<span class="chart-empty">График недоступен.</span>';
    return;
  }

  const xs = allPoints.map(p => p.x);
  const ys = allPoints.map(p => p.y);
  const minX = Math.min(...xs);
  const maxX = Math.max(...xs);
  const minY = Math.min(...ys);
  const maxY = Math.max(...ys);
  const marginX = (maxX - minX) * 0.08 || 1;
  const marginY = (maxY - minY) * 0.08 || 1;

  const width = 900;
  const height = 320;
  const padding = 48;
  const usableW = width - padding * 2;
  const usableH = height - padding * 2;
  const rangeX = (maxX - minX) + 2*marginX || 1;
  const rangeY = (maxY - minY) + 2*marginY || 1;

  const toX = (x) => padding + ((x - (minX - marginX)) / rangeX) * usableW;
  const toY = (y) => height - padding - ((y - (minY - marginY)) / rangeY) * usableH;

  const colors = ['#38bdf8', '#a78bfa', '#fb7185', '#34d399'];

  const polylines = seriesArray.map((s, idx) => {
    const pts = (s.points || []).map(p => `${toX(p.x)},${toY(p.y)}`).join(' ');
    return `<polyline fill="none" stroke="${colors[idx % colors.length]}" stroke-width="2" points="${pts}" />`;
  }).join('\n');

  const pointCircles = seriesArray.map((s, idx) => {
    return (s.points || []).map(p => `<circle cx="${toX(p.x)}" cy="${toY(p.y)}" r="3" fill="${colors[idx % colors.length]}" />`).join('');
  }).join('\n');

  const xTicks = Array.from({ length: 5 }, (_, i) => {
    const rx = i/4;
    const xv = (minX - marginX) + rx * ((maxX - minX) + 2*marginX);
    return { x: toX(xv), v: xv };
  });
  const yTicks = Array.from({ length: 5 }, (_, i) => {
    const ry = i/4;
    const yv = (minY - marginY) + (1-ry) * ((maxY - minY) + 2*marginY);
    return { y: toY(yv), v: yv };
  });

  const xTickLines = xTicks.map(t => `<text x="${t.x}" y="${height - 14}" text-anchor="middle" fill="#98a8c7" font-size="12">${t.v.toFixed(4)}</text>`).join('');
  const yTickLines = yTicks.map(t => `<g><line x1="${padding}" y1="${t.y}" x2="${width - padding}" y2="${t.y}" stroke="rgba(255,255,255,0.06)" /><text x="12" y="${t.y + 4}" fill="#98a8c7" font-size="12">${t.v.toFixed(4)}</text></g>`).join('');

  container.innerHTML = `
    <svg viewBox="0 0 ${width} ${height}" class="chart-svg" role="img">
      ${yTickLines}
      <line x1="${padding}" y1="${padding}" x2="${padding}" y2="${height - padding}" stroke="rgba(255,255,255,0.18)" />
      <line x1="${padding}" y1="${height - padding}" x2="${width - padding}" y2="${height - padding}" stroke="rgba(255,255,255,0.18)" />
      ${xTickLines}
      ${polylines}
      ${pointCircles}
    </svg>
  `;
}

function fillPreset8() {
  if (nodesCsv8Input) {
    nodesCsv8Input.value = formatTask7Table(task7PresetPoints);
  }
}

async function solveTask8() {
  try {
    const csv = nodesCsv8Input.value || '';
    const x1 = parseFloat(x1Input.value);
    const x2 = parseFloat(x2Input.value);
    const pts = [];
    if (!Number.isNaN(x1)) pts.push(x1);
    if (!Number.isNaN(x2)) pts.push(x2);

    const response = await fetch('/api/task8/solve', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ csv, evaluate_points: pts }),
    });
    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    const data = await response.json();
    renderDerivativesResults(data);
  } catch (err) {
    console.error(err);
    derivativesResult.innerHTML = `<span class="chart-empty">Ошибка: ${String(err)}</span>`;
  }
}

function renderDerivativesResults(response) {
  const nodes = response.nodes || [];
  const rows = response.evaluations || [];
  if (!Array.isArray(rows) || rows.length === 0) {
    derivativesResult.innerHTML = '<span class="chart-empty">Результаты отсутствуют.</span>';
  } else {
    let html = '<table class="data-table"><thead><tr><th>#</th><th>x</th><th>S(x)</th><th>S\'(x)</th><th>S\"(x)</th></tr></thead><tbody>';
    for (let i = 0; i < rows.length; i++) {
      const r = rows[i];
      html += `<tr><td>${i+1}</td><td>${Number(r.x).toFixed(6)}</td><td>${Number(r.s).toFixed(6)}</td><td>${Number(r.s1).toFixed(6)}</td><td>${Number(r.s2).toFixed(6)}</td></tr>`;
    }
    html += '</tbody></table>';
    derivativesResult.innerHTML = html;
  }
  renderNodesTable(nodesTable8Wrap, nodes);
  renderInterpolationChart(derivativesChartWrap, nodes);
}

async function solveTask9() {
  try {
    const equation = equation9Input.value || 'cos(x^2)/(x+1)';
    const a = parseFloat(a9Input.value) || 0;
    const b = parseFloat(b9Input.value) || 1;
    const eps = parseFloat(eps9Input.value) || 1e-4;

    const response = await fetch('/api/task9/solve', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ equation, a, b, epsilon: eps, max_iterations: 12 }),
    });
    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    const data = await response.json();
    // render iterations table using existing renderIterationTableIn
    const itRows = (data.iterations || []).map((row, idx) => ({ iteration: row.iteration, x: row.x, delta: row.delta, residual: row.residual }));
    renderIterationTableIn(iterationTable9Wrap, itRows);
    renderIterationChartIn(iterationChart9Wrap, itRows);
    integralResult.innerHTML = `<div>Интеграл ≈ ${Number(data.result).toFixed(8)}</div>`;
  } catch (err) {
    console.error(err);
    integralResult.innerHTML = `<span class="chart-empty">Ошибка: ${String(err)}</span>`;
    iterationTable9Wrap.innerHTML = '<span class="chart-empty">Таблица недоступна.</span>';
    iterationChart9Wrap.innerHTML = '<span class="chart-empty">График недоступен.</span>';
  }
}

async function solveTask10() {
  try {
    const equation = equation10Input.value || '1 + 0.8*y*sin(x) - 2*y*y';
    const a = parseFloat(a10Input.value) || 0;
    const b = parseFloat(b10Input.value) || 1;
    const y0 = parseFloat(y0Input.value) || 0;
    const h = parseFloat(h10Input.value) || 0.1;

    const response = await fetch('/api/task10/solve', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ equation, a, b, y0, h }),
    });
    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    const data = await response.json();

    // RK4 table rows
    const rkRows = (data.rk4 || []).map((item, idx) => ({ iteration: idx, x: item.x, delta: 0, residual: 0 }));
    renderIterationTableIn(rk4TableWrap, rkRows);
    // chart series
    const rkSeries = { label: 'RK4', points: (data.rk4 || []).map(p => ({ x: p.x, y: p.y })) };
    renderLineChart(rk4ChartWrap, [rkSeries]);
    rk4Result.innerHTML = `<div>Последняя точка: x=${Number((data.rk4 || []).slice(-1)[0]?.x || 0).toFixed(6)}, y=${Number((data.rk4 || []).slice(-1)[0]?.y || 0).toFixed(6)}</div>`;

    // Adams table
    const adRows = (data.adams || []).map((row, idx) => ({ iteration: idx, x: row.x, delta: row.delta || 0, residual: 0, y_pred: row.y_pred, y_corr: row.y_corr }));
    // transform to structure expected by renderIterationTableIn: it will show x and delta
    const adForTable = (data.adams || []).map((row, idx) => ({ iteration: idx, x: row.x, delta: row.delta || 0, residual: 0 }));
    renderIterationTableIn(adamsTableWrap, adForTable);
    const adSeries = { label: 'Adams', points: (data.adams || []).map(p => ({ x: p.x, y: p.y_corr || p.y_pred })) };
    renderLineChart(adamsChartWrap, [adSeries]);
    adamsResult.innerHTML = `<div>Последняя точка: x=${Number((data.adams || []).slice(-1)[0]?.x || 0).toFixed(6)}, y=${Number((data.adams || []).slice(-1)[0]?.y_corr || (data.adams || []).slice(-1)[0]?.y_pred || 0).toFixed(6)}</div>`;
  } catch (err) {
    console.error(err);
    rk4Result.innerHTML = `<span class="chart-empty">Ошибка: ${String(err)}</span>`;
    adamsResult.innerHTML = `<span class="chart-empty">Ошибка: ${String(err)}</span>`;
    rk4TableWrap.innerHTML = '<span class="chart-empty">Таблица недоступна.</span>';
    rk4ChartWrap.innerHTML = '<span class="chart-empty">График недоступен.</span>';
    adamsTableWrap.innerHTML = '<span class="chart-empty">Таблица недоступна.</span>';
    adamsChartWrap.innerHTML = '<span class="chart-empty">График недоступен.</span>';
  }
}

function renderTask11Table(rows) {
  if (!Array.isArray(rows) || rows.length === 0) {
    task11TableWrap.innerHTML = '<span class="chart-empty">Таблица недоступна.</span>';
    return;
  }
  const body = rows.map((r) => `
    <tr>
      <td>${r.index}</td>
      <td>${Number(r.x).toFixed(6)}</td>
      <td>${Number(r.a).toFixed(6)}</td>
      <td>${Number(r.b).toFixed(6)}</td>
      <td>${Number(r.c).toFixed(6)}</td>
      <td>${Number(r.d).toFixed(6)}</td>
      <td>${Number(r.alpha).toFixed(6)}</td>
      <td>${Number(r.beta).toFixed(6)}</td>
      <td>${Number(r.y).toFixed(6)}</td>
    </tr>
  `).join('');

  task11TableWrap.innerHTML = `
    <table class="data-table">
      <thead>
        <tr>
          <th>i</th>
          <th>x_i</th>
          <th>a_i</th>
          <th>b_i</th>
          <th>c_i</th>
          <th>d_i</th>
          <th>alpha_i</th>
          <th>beta_i</th>
          <th>y_i</th>
        </tr>
      </thead>
      <tbody>${body}</tbody>
    </table>
  `;
}

async function solveTask11() {
  try {
    const payload = {
      p_expr: (p11Input?.value || '2/x').trim(),
      q_expr: (q11Input?.value || '-3').trim(),
      rhs_expr: (rhs11Input?.value || '2').trim(),
      a: parseFloat(a11Input?.value || '0.8'),
      b: parseFloat(b11Input?.value || '1.1'),
      h: parseFloat(h11Input?.value || '0.1'),
      left_alpha: parseFloat(leftAlpha11Input?.value || '0'),
      left_beta: parseFloat(leftBeta11Input?.value || '1'),
      left_gamma: parseFloat(leftGamma11Input?.value || '1.5'),
      right_alpha: parseFloat(rightAlpha11Input?.value || '2'),
      right_beta: parseFloat(rightBeta11Input?.value || '1'),
      right_gamma: parseFloat(rightGamma11Input?.value || '3'),
    };

    const response = await fetch('/api/task11/solve', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    });
    if (!response.ok) {
      const text = await response.text();
      throw new Error(`HTTP ${response.status}${text ? `: ${text}` : ''}`);
    }
    const data = await response.json();

    const points = data.points || [];
    const rows = data.sweep || [];
    renderTask11Table(rows);
    renderLineChart(task11ChartWrap, [{ label: 'y(x)', points: points.map((p) => ({ x: p.x, y: p.y })) }]);

    const last = points.length ? points[points.length - 1] : null;
    task11Result.innerHTML = `
      <div>Узлов: ${points.length}</div>
      <div>Шаг (использован): ${Number(data.step_used || 0).toFixed(6)}</div>
      <div>Последняя точка: ${last ? `x=${Number(last.x).toFixed(6)}, y=${Number(last.y).toFixed(6)}` : 'нет'}</div>
    `;
  } catch (error) {
    console.error(error);
    task11Result.innerHTML = `<span class="chart-empty">Ошибка: ${String(error)}</span>`;
    task11TableWrap.innerHTML = '<span class="chart-empty">Таблица недоступна.</span>';
    task11ChartWrap.innerHTML = '<span class="chart-empty">График недоступен.</span>';
  }
}

async function solveTask7() {
  try {
    // read custom CSV if provided, otherwise use preset points
    let points = task7PresetPoints;
    try {
      const csv = nodesCsv7Input ? nodesCsv7Input.value.trim() : '';
      if (csv) {
        points = parseCsvPoints(csv);
      }
    } catch (err) {
      throw err;
    }

    const response = await fetch('/api/task7/solve', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ points }),
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }

    const data = await response.json();
    renderLeastSquaresResults(data);
  } catch (error) {
    console.error(error);
    linearResult.innerHTML = `<span class="chart-empty">Ошибка: ${String(error)}</span>`;
    quadraticResult.innerHTML = '<span class="chart-empty">Ошибка расчета.</span>';
    leastSquaresTableWrap.innerHTML = '<span class="chart-empty">Таблица недоступна.</span>';
    leastSquaresChartWrap.innerHTML = '<span class="chart-empty">График недоступен.</span>';
  }
}

// Клиентская Lagrange функция для построения графика
function evaluateLagrange(nodes, x) {
  const n = nodes.length;
  let res = 0;
  for (let i = 0; i < n; i++) {
    const xi = nodes[i].x;
    const yi = nodes[i].y;
    let li = 1;
    for (let j = 0; j < n; j++) {
      if (i === j) continue;
      const xj = nodes[j].x;
      li *= (x - xj) / (xi - xj);
    }
    res += yi * li;
  }
  return res;
}

function renderNodesTable(container, nodes) {
  if (!Array.isArray(nodes) || nodes.length === 0) {
    container.innerHTML = '<span class="chart-empty">Нет узлов.</span>';
    return;
  }
  const rows = nodes.map((n) => `<tr><td>${n.index}</td><td>${n.x.toFixed(6)}</td><td>${n.y.toFixed(6)}</td></tr>`).join('');
  container.innerHTML = `
    <table class="data-table">
      <thead><tr><th>#</th><th>x</th><th>y</th></tr></thead>
      <tbody>${rows}</tbody>
    </table>
  `;
}

function renderInterpolationChart(container, nodes) {
  if (!Array.isArray(nodes) || nodes.length === 0) {
    container.innerHTML = '<span class="chart-empty">График недоступен.</span>';
    return;
  }
  const xs = nodes.map((n) => n.x);
  const ys = nodes.map((n) => n.y);
  const minX = Math.min(...xs);
  const maxX = Math.max(...xs);
  const margin = (maxX - minX) * 0.1 || 1;
  const samples = 200;
  const pts = [];
  for (let i = 0; i <= samples; i++) {
    const x = minX - margin + (i / samples) * (maxX - minX + 2 * margin);
    const y = evaluateLagrange(nodes, x);
    pts.push({ x, y });
  }

  const width = 900;
  const height = 320;
  const padding = 48;
  const allY = pts.map((p) => p.y).concat(ys);
  const minY = Math.min(...allY);
  const maxY = Math.max(...allY);

  const usableW = width - padding * 2;
  const usableH = height - padding * 2;

  const xStart = minX - margin;
  const xRange = maxX - minX + 2 * margin || 1;
  const toX = (x) => padding + ((x - xStart) / xRange) * usableW;
  const toY = (y) => height - padding - ((y - minY) / (maxY - minY || 1)) * usableH;

  const polyline = pts.map((p) => `${toX(p.x)},${toY(p.y)}`).join(' ');

  const nodesCircles = nodes.map((n) => `<circle cx="${toX(n.x)}" cy="${toY(n.y)}" r="4" fill="#f97316"><title>(${n.x.toFixed(6)}, ${n.y.toFixed(6)})</title></circle>`).join('');

  // axis ticks (5)
  const xTicks = Array.from({ length: 5 }, (_, i) => {
    const rx = i / 4;
    const xv = xStart + rx * xRange;
    return { x: toX(xv), v: xv };
  });
  const yTicks = Array.from({ length: 5 }, (_, i) => {
    const ry = i / 4;
    const yv = minY + (1 - ry) * (maxY - minY);
    return { y: toY(yv), v: yv };
  });

  const xTickLines = xTicks.map((t) => `<text x="${t.x}" y="${height - 14}" text-anchor="middle" fill="#98a8c7" font-size="12">${t.v.toFixed(4)}</text>`).join('');
  const yTickLines = yTicks.map((t) => `<g><line x1="${padding}" y1="${t.y}" x2="${width - padding}" y2="${t.y}" stroke="rgba(255,255,255,0.06)" /><text x="12" y="${t.y + 4}" fill="#98a8c7" font-size="12">${t.v.toFixed(4)}</text></g>`).join('');

  container.innerHTML = `
    <svg viewBox="0 0 ${width} ${height}" class="chart-svg">
      ${yTickLines}
      <line x1="${padding}" y1="${padding}" x2="${padding}" y2="${height - padding}" stroke="rgba(255,255,255,0.18)" />
      <line x1="${padding}" y1="${height - padding}" x2="${width - padding}" y2="${height - padding}" stroke="rgba(255,255,255,0.18)" />
      ${xTickLines}
      <polyline fill="none" stroke="#60a5fa" stroke-width="2" points="${polyline}" />
      ${nodesCircles}
    </svg>
  `;
}

async function solveTask5() {
  const csv = nodesCsvInput.value || '';
  const xStar = parseFloat(xStarInput.value) || 0;
  try {
    const body = new URLSearchParams();
    body.set('csv', csv);
    body.set('x_star', String(xStar));
    const response = await fetch('/api/task5/solve', {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: body.toString(),
    });
    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    const data = await response.json();
    renderTask5Response(data);
  } catch (err) {
    console.error(err);
    lagrangeResult.innerHTML = `<span class="chart-empty">Ошибка: ${String(err)}</span>`;
    newtonResult.innerHTML = `<span class="chart-empty">Ошибка: ${String(err)}</span>`;
    nodesTableWrap.innerHTML = '<span class="chart-empty">Таблица недоступна.</span>';
    interpolationChartWrap.innerHTML = '<span class="chart-empty">График недоступен.</span>';
  }
}

function renderTask5Response(data) {
  const nodes = data.nodes || [];
  lagrangeResult.innerHTML = `<div>Значение (Лагранж) = ${Number(data.lagrange).toFixed(6)}<br>Оценка погрешности ≈ ${Number(data.err_lagrange).toExponential(3)}</div>`;
  newtonResult.innerHTML = `<div>Значение (Ньютон) = ${Number(data.newton).toFixed(6)}<br>Оценка погрешности ≈ ${Number(data.err_newton).toExponential(3)}</div>`;
  renderNodesTable(nodesTableWrap, nodes);
  renderInterpolationChart(interpolationChartWrap, nodes);
}

async function solveTask4Equation(payload, target) {
  try {
    const response = await fetch('/api/task4/solve', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }

    const data = await response.json();
    renderTask4Response(data, target);
  } catch (error) {
    console.error(error);
    target.bisectionResult.innerHTML = `<span class="chart-empty">Ошибка: ${String(error)}</span>`;
    target.simpleIterationResult.innerHTML = '<span class="chart-empty">Ошибка расчета.</span>';
    target.newtonResult.innerHTML = '<span class="chart-empty">Ошибка расчета.</span>';
    target.bisectionTableWrap.innerHTML = '<span class="chart-empty">Таблица недоступна.</span>';
    target.bisectionChartWrap.innerHTML = '<span class="chart-empty">График недоступен.</span>';
    target.simpleIterationTableWrap.innerHTML = '<span class="chart-empty">Таблица недоступна.</span>';
    target.simpleIterationChartWrap.innerHTML = '<span class="chart-empty">График недоступен.</span>';
    target.newtonTableWrap.innerHTML = '<span class="chart-empty">Таблица недоступна.</span>';
    target.newtonChartWrap.innerHTML = '<span class="chart-empty">График недоступен.</span>';
  }
}

function renderTask4Response(response, target) {
  const bisection = response.bisection || {};
  const simpleIteration = response.simple_iteration || {};
  const newton = response.newton || {};
  
  // Бисекция
  const bisectionLines = [];
  if (bisection.error) {
    bisectionLines.push(bisection.error);
  }
  bisectionLines.push(bisection.converged ? 'Метод сошёлся.' : 'Метод не сошёлся.');
  if (typeof bisection.final_delta === 'number') {
    bisectionLines.push(`Последняя δ = ${bisection.final_delta.toFixed(6)}`);
  }
  if (typeof bisection.final_residual === 'number') {
    bisectionLines.push(`Последняя невязка |f(x)| = ${bisection.final_residual.toFixed(6)}`);
  }
  if (bisection.solution !== undefined) {
    bisectionLines.push(`Решение: x = ${bisection.solution.toFixed(6)}`);
  }
  renderResultList(target.bisectionResult, bisectionLines, 'Бисекция');
  
  // Простая итерация
  const simpleIterationLines = [];
  if (simpleIteration.error) {
    simpleIterationLines.push(simpleIteration.error);
  }
  simpleIterationLines.push(simpleIteration.converged ? 'Метод сошёлся.' : 'Метод не сошёлся.');
  if (typeof simpleIteration.final_delta === 'number') {
    simpleIterationLines.push(`Последняя δ = ${simpleIteration.final_delta.toFixed(6)}`);
  }
  if (typeof simpleIteration.final_residual === 'number') {
    simpleIterationLines.push(`Последняя невязка |f(x)| = ${simpleIteration.final_residual.toFixed(6)}`);
  }
  if (simpleIteration.solution !== undefined) {
    simpleIterationLines.push(`Решение: x = ${simpleIteration.solution.toFixed(6)}`);
  }
  renderResultList(target.simpleIterationResult, simpleIterationLines, 'Простая итерация');
  
  // Метод Ньютона
  const newtonLines = [];
  if (newton.error) {
    newtonLines.push(newton.error);
  }
  newtonLines.push(newton.converged ? 'Метод сошёлся.' : 'Метод не сошёлся.');
  if (typeof newton.final_delta === 'number') {
    newtonLines.push(`Последняя δ = ${newton.final_delta.toFixed(6)}`);
  }
  if (typeof newton.final_residual === 'number') {
    newtonLines.push(`Последняя невязка |f(x)| = ${newton.final_residual.toFixed(6)}`);
  }
  if (newton.solution !== undefined) {
    newtonLines.push(`Решение: x = ${newton.solution.toFixed(6)}`);
  }
  renderResultList(target.newtonResult, newtonLines, 'Метод Ньютона');
  
  // Таблицы и графики
  renderIterationTableIn(target.bisectionTableWrap, bisection.iterations || []);
  renderIterationChartIn(target.bisectionChartWrap, bisection.iterations || []);
  renderIterationTableIn(target.simpleIterationTableWrap, simpleIteration.iterations || []);
  renderIterationChartIn(target.simpleIterationChartWrap, simpleIteration.iterations || []);
  renderIterationTableIn(target.newtonTableWrap, newton.iterations || []);
  renderIterationChartIn(target.newtonChartWrap, newton.iterations || []);
}

async function solvePreset4() {
  await solveTask4Equation({
    equation: presetEquation4.equation,
    a: presetEquation4.a,
    b: presetEquation4.b,
    x0: presetEquation4.x0,
    epsilon: presetEquation4.epsilon,
    max_iterations: presetEquation4.maxIterations,
  }, {
    bisectionResult: bisectionPresetResult,
    simpleIterationResult: simpleIterationPresetResult,
    newtonResult: newtonPresetResult,
    bisectionTableWrap: bisectionPresetIterationTableWrap,
    bisectionChartWrap: bisectionPresetIterationChartWrap,
    simpleIterationTableWrap: simpleIterationPresetIterationTableWrap,
    simpleIterationChartWrap: simpleIterationPresetIterationChartWrap,
    newtonTableWrap: newtonPresetIterationTableWrap,
    newtonChartWrap: newtonPresetIterationChartWrap,
  });
}

async function solveCustom4() {
  const payload = {
    equation: equationInput.value,
    a: parseFloat(equation4AInput.value),
    b: parseFloat(equation4BInput.value),
    x0: parseFloat(equation4X0Input.value),
    epsilon: parseFloat(equation4EpsilonInput.value),
    max_iterations: parseInt(equation4MaxIterationsInput.value, 10),
  };
  
  await solveTask4Equation(payload, {
    bisectionResult: bisectionCustomResult,
    simpleIterationResult: simpleIterationCustomResult,
    newtonResult: newtonCustomResult,
    bisectionTableWrap: bisectionCustomIterationTableWrap,
    bisectionChartWrap: bisectionCustomIterationChartWrap,
    simpleIterationTableWrap: simpleIterationCustomIterationTableWrap,
    simpleIterationChartWrap: simpleIterationCustomIterationChartWrap,
    newtonTableWrap: newtonCustomIterationTableWrap,
    newtonChartWrap: newtonCustomIterationChartWrap,
  });
}

// Обработчик построения произвольной СЛАУ.
buildMatrixButton.addEventListener('click', () => {
  buildMatrixEditor(systemSizeInput.value);
});

// Синхронизируем форму с выбранной размерностью при загрузке и изменении значения.
systemSizeInput.addEventListener('change', () => {
  buildMatrixEditor(systemSizeInput.value);
});

// Заполняем форму вариантом 13.
fillPresetButton.addEventListener('click', () => {
  fillPresetSystem();
});

// Сразу считаем вариант 13 по кнопке из верхней панели.
solvePresetButton.addEventListener('click', () => {
  solvePreset();
});

// Кнопка рядом с исходной системой быстро ведет к графику.
showGraphButton.addEventListener('click', () => {
  showGraph();
});

// Отправка произвольной СЛАУ из формы.
customSystemForm.addEventListener('submit', (event) => {
  event.preventDefault();
  solveCustom1();
});

buildMatrix2Button.addEventListener('click', () => {
  buildMatrixEditor2(system2SizeInput.value);
});

system2SizeInput.addEventListener('change', () => {
  buildMatrixEditor2(system2SizeInput.value);
});

fillPreset2Button.addEventListener('click', () => {
  fillPresetSystem2();
});

solvePreset2Button.addEventListener('click', () => {
  solvePreset2();
});

showGraph2Button.addEventListener('click', () => {
  showGraph2();
});

customSystem2Form.addEventListener('submit', (event) => {
  event.preventDefault();
  solveCustom2();
});

if (fillPreset3Button) {
  fillPreset3Button.addEventListener('click', () => {
    fillPreset3Info();
  });
}

if (saveTask3InfoButton) {
  saveTask3InfoButton.addEventListener('click', () => {
    saveTask3Info();
  });
}

  if (fillPreset4Button) {
    fillPreset4Button.addEventListener('click', () => {
      fillPresetEquation4();
    });
  }

  if (solvePreset4Button) {
    solvePreset4Button.addEventListener('click', () => {
      solvePreset4();
    });
  }

  if (customEquation4Form) {
    customEquation4Form.addEventListener('submit', (event) => {
      event.preventDefault();
      solveCustom4();
    });
  }

  if (fillPreset5Button) {
    fillPreset5Button.addEventListener('click', () => { fillPreset5(); });
  }

  if (solvePreset5Button) {
    solvePreset5Button.addEventListener('click', () => { solveTask5(); });
  }

  if (fillPreset6Button) {
    fillPreset6Button.addEventListener('click', () => {
      fillPreset6();
    });
  }

  if (solvePreset6Button) {
    solvePreset6Button.addEventListener('click', () => {
      solveTask6();
    });
  }

  if (fillPreset7Button) {
    fillPreset7Button.addEventListener('click', () => {
      fillPreset7();
    });
  }

  if (fillPreset8Button) {
    fillPreset8Button.addEventListener('click', () => {
      fillPreset8();
    });
  }

  if (fillPreset9Button) {
    fillPreset9Button.addEventListener('click', () => {
      if (equation9Input) equation9Input.value = 'cos(x^2)/(x+1)';
      if (a9Input) a9Input.value = '0';
      if (b9Input) b9Input.value = '1';
      if (eps9Input) eps9Input.value = '1e-4';
    });
  }

  if (fillPreset10Button) {
    fillPreset10Button.addEventListener('click', () => {
      if (equation10Input) equation10Input.value = "1 + 0.8*y*sin(x) - 2*y*y";
      if (a10Input) a10Input.value = '0';
      if (b10Input) b10Input.value = '1';
      if (y0Input) y0Input.value = '0';
      if (h10Input) h10Input.value = '0.1';
    });
  }

  if (fillPreset11Button) {
    fillPreset11Button.addEventListener('click', () => {
      if (p11Input) p11Input.value = '2/x';
      if (q11Input) q11Input.value = '-3';
      if (rhs11Input) rhs11Input.value = '2';
      if (a11Input) a11Input.value = '0.8';
      if (b11Input) b11Input.value = '1.1';
      if (h11Input) h11Input.value = '0.1';
      if (leftAlpha11Input) leftAlpha11Input.value = '0';
      if (leftBeta11Input) leftBeta11Input.value = '1';
      if (leftGamma11Input) leftGamma11Input.value = '1.5';
      if (rightAlpha11Input) rightAlpha11Input.value = '2';
      if (rightBeta11Input) rightBeta11Input.value = '1';
      if (rightGamma11Input) rightGamma11Input.value = '3';
    });
  }

  if (solvePreset7Button) {
    solvePreset7Button.addEventListener('click', () => {
      solveTask7();
    });
  }


  if (solvePreset8Button) {
    solvePreset8Button.addEventListener('click', () => {
      solveTask8();
    });
  }

  if (solvePreset9Button) {
    solvePreset9Button.addEventListener('click', () => {
      solveTask9();
    });
  }

  if (solvePreset10Button) {
    solvePreset10Button.addEventListener('click', () => {
      solveTask10();
    });
  }

  if (solvePreset11Button) {
    solvePreset11Button.addEventListener('click', () => {
      solveTask11();
    });
  }

// Инициализация интерфейса первого задания.
buildMatrixEditor(presetSystem.matrix.length, presetSystem);
fillPresetSystem();
buildMatrixEditor2(presetSystem2.matrix.length, presetSystem2);
fillPresetSystem2();

fillPresetEquation4();
fillPreset3Info();
fillPreset6();
fillPreset7();
fillPreset8();
if (fillPreset11Button) {
  fillPreset11Button.click();
}

loadTasks();
