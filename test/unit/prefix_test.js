// Prefix Unit Tests for HTML-Go Converter
const fetch = require('node-fetch');

// 设置后端服务器的基础URL
const BASE_URL = process.env.SERVER_URL || 'http://localhost:8080';
console.log(`Using server at: ${BASE_URL}`);

describe('Component Prefix Tests', () => {
  // 基础HTML测试 - 确保没有前缀问题
  test('Basic HTML conversion has correct package prefix', async () => {
    const testData = {
      html: '<div><h1>Test</h1></div>',
      packagePrefix: 'h',
      direction: 'html2go'
    };

    const response = await fetch(`${BASE_URL}/convert`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(testData),
    });

    expect(response.ok).toBeTruthy();
    const data = await response.json();

    expect(data.code).toContain('h.Div');
    expect(data.code).toContain('h.H1');
    // 确保没有出现错误的前缀
    expect(data.code).not.toMatch(/hh\./);
  });

  // Vuetify前缀测试 - 确保前缀是'v'
  test('Vuetify components use "v" prefix', async () => {
    const testData = {
      html: `<div>
        <v-btn>Test Button</v-btn>
        <v-card><v-card-title>Card</v-card-title></v-card>
      </div>`,
      packagePrefix: 'h',
      vuetifyPrefix: 'v',
      vuetifyXPrefix: 'vx',
      direction: 'html2go'
    };

    const response = await fetch(`${BASE_URL}/convert`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(testData),
    });

    expect(response.ok).toBeTruthy();
    const data = await response.json();

    // 检查Vuetify组件是否正确使用'v'前缀
    expect(data.code).toMatch(/v\.VBtn\(/);
    expect(data.code).toMatch(/v\.VCard\(/);
    expect(data.code).toMatch(/v\.VCardTitle\(/);

    // 确保没有错误的前缀组合
    expect(data.code).not.toMatch(/hv\./);

    console.log('Generated Vuetify code:', data.code);
  });

  // VuetifyX前缀测试 - 确保前缀是'vx'
  test('VuetifyX components use "vx" prefix', async () => {
    const testData = {
      html: `<div>
        <vx-dialog title="Test">Dialog Content</vx-dialog>
        <vx-date-picker></vx-date-picker>
      </div>`,
      packagePrefix: 'h',
      vuetifyPrefix: 'v',
      vuetifyXPrefix: 'vx',
      direction: 'html2go'
    };

    const response = await fetch(`${BASE_URL}/convert`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(testData),
    });

    expect(response.ok).toBeTruthy();
    const data = await response.json();

    // 检查VuetifyX组件是否正确使用'vx'前缀
    expect(data.code).toMatch(/vx\.VXDialog\(/);
    expect(data.code).toMatch(/vx\.VXDatepicker\(/);

    // 确保没有错误的前缀组合
    expect(data.code).not.toMatch(/hvx\./);

    console.log('Generated VuetifyX code:', data.code);
  });

  // 混合组件测试 - 同时测试Vuetify和VuetifyX组件
  test('Mixed components use correct prefixes', async () => {
    const testData = {
      html: `<div>
        <v-btn color="primary">Vuetify Button</v-btn>
        <vx-dialog title="Test Dialog">
          <v-card>
            <v-card-title>Card in Dialog</v-card-title>
          </v-card>
        </vx-dialog>
      </div>`,
      packagePrefix: 'h',
      vuetifyPrefix: 'v',
      vuetifyXPrefix: 'vx',
      direction: 'html2go'
    };

    const response = await fetch(`${BASE_URL}/convert`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(testData),
    });

    expect(response.ok).toBeTruthy();
    const data = await response.json();

    // 检查Vuetify组件
    expect(data.code).toMatch(/v\.VBtn\(/);
    expect(data.code).toMatch(/v\.VCard\(/);
    expect(data.code).toMatch(/v\.VCardTitle\(/);

    // 检查VuetifyX组件
    expect(data.code).toMatch(/vx\.VXDialog\(/);

    // 确保没有错误的前缀组合
    expect(data.code).not.toMatch(/hv\./);
    expect(data.code).not.toMatch(/hvx\./);

    console.log('Generated mixed components code:', data.code);
  });
}); 