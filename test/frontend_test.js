// Frontend API Tests for HTML-Go Converter
// This file tests the frontend API interactions with the backend service

// 使用node-fetch替代全局fetch来进行实际网络请求
const fetch = require('node-fetch');
// 设置后端服务器的基础URL
const BASE_URL = 'http://localhost:8080';

describe('Frontend API Tests', () => {
  // 测试基本HTML转Go代码功能
  test('HTML to Go basic conversion succeeds', async () => {
    // 测试数据
    const testData = {
      html: '<div class="container"><h1 class="text-xl font-bold">Hello World</h1></div>',
      packagePrefix: 'h',
      vuetifyPrefix: 'v',
      vuetifyXPrefix: 'vx',
      childrenMode: false,
      direction: 'html2go'
    };

    // 发起实际请求到后端服务
    const response = await fetch(`${BASE_URL}/convert`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(testData),
    });

    // 验证响应
    expect(response.ok).toBeTruthy();
    expect(response.status).toBe(200);

    const data = await response.json();
    expect(data.code).toBeDefined();
    expect(data.error).toBeFalsy();

    // 验证生成的代码包含预期的元素
    expect(data.code).toContain('h.Div');
    expect(data.code).toContain('h.H1');
    expect(data.code).toContain('Class("container")');
    expect(data.code).toContain('Class("text-xl font-bold")');
  });

  // 测试Vuetify组件转换
  test('Vuetify component conversion succeeds with custom prefix', async () => {
    // 测试数据 - 使用自定义前缀
    const customVuetifyPrefix = 'custom_vuetify';
    const testData = {
      html: `<div>
        <v-btn color="primary">Click me</v-btn>
        <v-card>
          <v-card-title>Card Title</v-card-title>
          <v-card-text>Card content goes here</v-card-text>
        </v-card>
      </div>`,
      packagePrefix: 'h',
      vuetifyPrefix: customVuetifyPrefix,
      vuetifyXPrefix: 'vx',
      childrenMode: false,
      direction: 'html2go'
    };

    // 发起实际请求到后端服务
    const response = await fetch(`${BASE_URL}/convert`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(testData),
    });

    // 验证响应
    expect(response.ok).toBeTruthy();
    expect(response.status).toBe(200);

    const data = await response.json();
    expect(data.code).toBeDefined();
    expect(data.error).toBeFalsy();

    // 验证生成的代码包含预期的元素
    expect(data.code).toContain(`${customVuetifyPrefix}.VBtn`);
    expect(data.code).toContain(`${customVuetifyPrefix}.VCard`);
    expect(data.code).toContain(`${customVuetifyPrefix}.VCardTitle`);
    expect(data.code).toContain(`${customVuetifyPrefix}.VCardText`);
    expect(data.code).toContain('Color("primary")');
    expect(data.code).toContain('Card Title');
    expect(data.code).toContain('Card content goes here');
  });

  // 测试VuetifyX组件转换
  test('VuetifyX component conversion succeeds with custom prefix', async () => {
    // 测试数据 - 使用自定义前缀
    const customVuetifyXPrefix = 'custom_vuetifyx';
    const testData = {
      html: `<div>
        <vx-date-picker label="Select Date"></vx-date-picker>
        <vx-dialog title="Confirmation" text="Are you sure?">
          <v-btn color="primary">Open Dialog</v-btn>
        </vx-dialog>
      </div>`,
      packagePrefix: 'h',
      vuetifyPrefix: 'v',
      vuetifyXPrefix: customVuetifyXPrefix,
      childrenMode: false,
      direction: 'html2go'
    };

    // 发起实际请求到后端服务
    const response = await fetch(`${BASE_URL}/convert`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(testData),
    });

    // 验证响应
    expect(response.ok).toBeTruthy();
    expect(response.status).toBe(200);

    const data = await response.json();
    expect(data.code).toBeDefined();
    expect(data.error).toBeFalsy();

    // 验证生成的代码包含预期的元素
    expect(data.code).toContain(`${customVuetifyXPrefix}.VXDatepicker`);
    expect(data.code).toContain(`${customVuetifyXPrefix}.VXDialog`);
    expect(data.code).toContain('Attr("label", "Select Date")');
    expect(data.code).toContain('Title("Confirmation")');
    expect(data.code).toContain('v.VBtn');
    expect(data.code).toContain('Color("primary")');
  });

  // 测试默认前缀设置
  test('Default prefixes are handled correctly', async () => {
    // 测试数据 - 不设置前缀
    const testData = {
      html: `<div>
        <v-btn color="primary">Click me</v-btn>
        <vx-date-picker label="Select Date"></vx-date-picker>
      </div>`,
      direction: 'html2go'
    };

    // 发起实际请求到后端服务
    const response = await fetch(`${BASE_URL}/convert`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(testData),
    });

    // 验证响应
    expect(response.ok).toBeTruthy();
    expect(response.status).toBe(200);

    const data = await response.json();
    expect(data.code).toBeDefined();
    expect(data.error).toBeFalsy();

    // 验证生成的代码包含预期的元素 (使用默认前缀)
    expect(data.code).toContain('v.VBtn');
    expect(data.code).toContain('vx.VXDatepicker');
    expect(data.code).toContain('Color("primary")');
    expect(data.code).toContain('Attr("label", "Select Date")');
  });

  // 专门测试Vuetify和VuetifyX默认前缀
  test('Vuetify and VuetifyX default prefixes are correctly set to "v" and "vx"', async () => {
    // 测试数据 - 使用默认前缀
    const testData = {
      html: `<div>
        <v-btn>Vuetify Button</v-btn>
        <v-card>
          <v-card-title>Vuetify Card</v-card-title>
        </v-card>
        <vx-dialog title="VuetifyX Dialog">Dialog Content</vx-dialog>
        <vx-date-picker></vx-date-picker>
      </div>`,
      packagePrefix: 'h',
      // 明确设置为v和vx
      vuetifyPrefix: 'v',
      vuetifyXPrefix: 'vx',
      direction: 'html2go'
    };

    // 发起实际请求到后端服务
    const response = await fetch(`${BASE_URL}/convert`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(testData),
    });

    // 验证响应
    expect(response.ok).toBeTruthy();
    expect(response.status).toBe(200);

    const data = await response.json();
    expect(data.code).toBeDefined();
    expect(data.error).toBeFalsy();

    // 严格检查生成的代码是否包含以v.V开头的Vuetify组件
    expect(data.code).toMatch(/v\.VBtn\(/);
    expect(data.code).toMatch(/v\.VCard\(/);
    expect(data.code).toMatch(/v\.VCardTitle\(/);

    // 严格检查生成的代码是否包含以vx.VX开头的VuetifyX组件
    expect(data.code).toMatch(/vx\.VXDialog\(/);
    expect(data.code).toMatch(/vx\.VXDatepicker\(/);

    // 确保没有错误的前缀如hv或hvx
    expect(data.code).not.toMatch(/hv\./);
    expect(data.code).not.toMatch(/hvx\./);
  });

  // 测试复杂嵌套HTML结构
  test('Complex nested HTML conversion succeeds', async () => {
    // 测试数据 - 复杂嵌套HTML
    const customPackagePrefix = 'custom_pkg';
    const testData = {
      html: `<div class="container">
        <header class="header">
          <nav class="navbar">
            <ul class="nav-list">
              <li class="nav-item"><a href="#" class="nav-link">Home</a></li>
              <li class="nav-item"><a href="#" class="nav-link">About</a></li>
              <li class="nav-item"><a href="#" class="nav-link">Contact</a></li>
            </ul>
          </nav>
        </header>
        <main class="content">
          <section class="hero">
            <h1 class="title">Welcome to our site</h1>
            <p class="subtitle">This is a complex HTML structure</p>
          </section>
        </main>
      </div>`,
      packagePrefix: customPackagePrefix,
      direction: 'html2go'
    };

    // 发起实际请求到后端服务
    const response = await fetch(`${BASE_URL}/convert`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(testData),
    });

    // 验证响应
    expect(response.ok).toBeTruthy();
    expect(response.status).toBe(200);

    const data = await response.json();
    expect(data.code).toBeDefined();
    expect(data.error).toBeFalsy();

    // 验证生成的代码包含预期的元素
    expect(data.code).toContain(`${customPackagePrefix}.Div`);
    expect(data.code).toContain(`${customPackagePrefix}.Header`);
    expect(data.code).toContain(`${customPackagePrefix}.Nav`);
    expect(data.code).toContain(`${customPackagePrefix}.Ul`);
    expect(data.code).toContain(`${customPackagePrefix}.Li`);
    expect(data.code).toContain(`${customPackagePrefix}.H1`);
    expect(data.code).toContain(`${customPackagePrefix}.P`);
    expect(data.code).toContain('Class("container")');
    expect(data.code).toContain('Class("nav-link")');
  });

  // 测试Children模式
  test('Children mode conversion succeeds', async () => {
    // 测试数据 - 使用Children模式
    const testData = {
      html: `<div><span>Text</span><v-btn>Click</v-btn></div>`,
      packagePrefix: 'h',
      vuetifyPrefix: 'v',
      vuetifyXPrefix: 'vx',
      childrenMode: true,
      direction: 'html2go'
    };

    // 发起实际请求到后端服务
    const response = await fetch(`${BASE_URL}/convert`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(testData),
    });

    // 验证响应
    expect(response.ok).toBeTruthy();
    expect(response.status).toBe(200);

    const data = await response.json();
    expect(data.code).toBeDefined();
    expect(data.error).toBeFalsy();

    // 验证生成的代码包含预期的元素
    expect(data.code).toContain('Children(');
    expect(data.code).toContain('h.Span("Text")');
    expect(data.code).toContain('v.VBtn');
    expect(data.code).toContain('"Click"');
  });

  // 测试语法修复功能
  test('Syntax fixing works correctly', async () => {
    // 测试数据 - 生成可能需要语法修复的代码
    const testData = {
      html: `<div id="app" class="container"></div>`,
      direction: 'html2go'
    };

    // 发起实际请求到后端服务
    const response = await fetch(`${BASE_URL}/convert`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(testData),
    });

    // 验证响应
    expect(response.ok).toBeTruthy();
    expect(response.status).toBe(200);

    const data = await response.json();
    expect(data.code).toBeDefined();
    expect(data.error).toBeFalsy();

    // 验证语法修复效果 - 确认基本元素存在
    expect(data.code).toContain('Div(');
    expect(data.code).toContain('Id("app")');
    expect(data.code).toContain('Class("container")');

    // 确保没有明显的语法错误
    expect(data.code).not.toContain('.Class');  // 不应该有前缀点号
    expect(data.code).not.toContain(')Class');  // 确保是 ).Class 而不是 )Class

    // 检查代码是否符合有效的Go语法 (不检查特定格式)
    // 确保左右括号数量相等
    const openParens = (data.code.match(/\(/g) || []).length;
    const closeParens = (data.code.match(/\)/g) || []).length;
    expect(openParens).toEqual(closeParens);
  });

  // 测试错误处理
  test('Various invalid requests are handled correctly', async () => {
    // 测试空HTML输入
    const emptyHTMLTest = {
      html: '',
      packagePrefix: 'h',
      direction: 'html2go'
    };

    let response = await fetch(`${BASE_URL}/convert`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(emptyHTMLTest),
    });

    // 验证错误响应
    expect(response.ok).toBeFalsy();
    expect(response.status).toBe(400);
    let data = await response.json();
    expect(data.error).toContain('HTML content is required');

    // 测试无效转换方向
    const invalidDirectionTest = {
      html: '<div>Test</div>',
      packagePrefix: 'h',
      direction: 'invalid_direction'
    };

    response = await fetch(`${BASE_URL}/convert`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(invalidDirectionTest),
    });

    // 验证错误响应
    expect(response.ok).toBeFalsy();
    expect(response.status).toBe(400);
    data = await response.json();
    expect(data.error).toContain('Invalid conversion direction');

    // 测试未实现的Go到HTML转换
    const go2htmlTest = {
      html: '<div>Test</div>',
      packagePrefix: 'h',
      direction: 'go2html'
    };

    response = await fetch(`${BASE_URL}/convert`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(go2htmlTest),
    });

    // 验证错误响应
    expect(response.ok).toBeFalsy();
    expect(response.status).toBe(501); // Not Implemented
    data = await response.json();
    expect(data.error).toContain('Go to HTML conversion is not implemented yet');
  });

  // 测试HTTP方法验证
  test('Only POST method is accepted', async () => {
    // 测试GET方法
    const response = await fetch(`${BASE_URL}/convert`, {
      method: 'GET',
    });

    // 验证错误响应
    expect(response.ok).toBeFalsy();
    expect(response.status).toBe(405); // Method Not Allowed
  });
});