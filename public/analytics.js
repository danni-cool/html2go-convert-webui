// Initialize Vercel Analytics
import { inject } from 'https://esm.sh/@vercel/analytics';

// Inject Vercel Analytics
inject();

// App version
const APP_VERSION = '1.0.0';

// Determine the current environment
const determineEnvironment = function () {
  // Check for manual override via URL parameter (?env=prod or ?env=local)
  const urlParams = new URLSearchParams(window.location.search);
  const envParam = urlParams.get('env');
  if (envParam === 'prod' || envParam === 'local') {
    return envParam;
  }

  const hostname = window.location.hostname;
  const port = window.location.port;

  // Check if URL contains development indicators
  const isDev = hostname === 'localhost' ||
    hostname === '127.0.0.1' ||
    hostname.startsWith('192.168.') ||
    hostname.startsWith('10.') ||
    port === '3000' ||
    port === '5000' ||
    port === '8000' ||
    port === '8080' ||
    hostname.includes('dev.') ||
    hostname.includes('-dev') ||
    (hostname.includes('vercel.app') && (hostname.includes('preview') || hostname.includes('dev-')));

  return isDev ? 'local' : 'prod';
};

// Get the current environment
const currentEnvironment = determineEnvironment();
console.log(`Analytics environment: ${currentEnvironment}`);

// Create a global Analytics object to handle all tracking
window.Analytics = {
  // Current environment
  environment: currentEnvironment,

  // App version
  version: APP_VERSION,

  // Prepare event data with environment info
  prepareEventData: function (data = {}) {
    return {
      ...data,
      environment: this.environment,
      version: this.version
    };
  },

  // Track application initialization
  // trackAppInit: function () {
  //   if (window.va) {
  //     window.va.track('app_initialized', this.prepareEventData({
  //       screen_width: window.innerWidth,
  //       screen_height: window.innerHeight,
  //       user_agent: navigator.userAgent,
  //       language: navigator.language,
  //       timestamp: new Date().toISOString()
  //     }));
  //   }
  // },

  // Track editor creation
  // trackEditorCreated: function (type) {
  //   if (window.va) {
  //     window.va.track('editor_created', this.prepareEventData({ type: type }));
  //   }
  // },

  // Track editor focus
  // trackEditorFocus: function (editorType) {
  //   if (window.va) {
  //     window.va.track(`focus_${editorType}_editor`, this.prepareEventData());
  //   }
  // },

  // Track editor content changes
  // trackEditorEdit: function (editorType, data) {
  //   if (window.va) {
  //     window.va.track(`edit_${editorType}`, this.prepareEventData(data));
  //   }
  // },

  // Track text selection
  // trackTextSelection: function (editorType, data) {
  //   if (window.va) {
  //     window.va.track(`select_${editorType}_text`, this.prepareEventData(data));
  //   }
  // },

  // Track button clicks
  // trackButtonClick: function (buttonName, data = {}) {
  //   if (window.va) {
  //     window.va.track(buttonName, this.prepareEventData(data));
  //   }
  // },

  // Utility function for creating debounced tracking functions
  createDebouncedTracker: function (editorInstance, editorType) {
    // Content change tracking with debounce
    let editDebounceTimer;
    const trackEditWithDebounce = function () {
      clearTimeout(editDebounceTimer);
      editDebounceTimer = setTimeout(function () {
        if (editorInstance) {
          window.Analytics.trackEditorEdit(editorType, {
            length: editorInstance.getValue().length,
            lines: editorInstance.getModel().getLineCount()
          });
        }
      }, 1500);
    };

    // Selection tracking with debounce
    let selectionTimer;
    const trackSelectionWithDebounce = function (e) {
      clearTimeout(selectionTimer);
      selectionTimer = setTimeout(function () {
        if (editorInstance) {
          const selection = editorInstance.getSelection();
          if (selection && !selection.isEmpty()) {
            window.Analytics.trackTextSelection(editorType, {
              startLineNumber: selection.startLineNumber,
              endLineNumber: selection.endLineNumber,
              length: editorInstance.getModel().getValueInRange(selection).length
            });
          }
        }
      }, 1000);
    };

    return {
      trackEdit: trackEditWithDebounce,
      trackSelection: trackSelectionWithDebounce,
      trackFocus: function () {
        window.Analytics.trackEditorFocus(editorType);
      }
    };
  },

  // Setup all editor analytics
  setupEditorAnalytics: function (htmlEditor, goEditor) {
    if (htmlEditor) {
      const htmlTracker = this.createDebouncedTracker(htmlEditor, 'html');

      // Setup event listeners
      htmlEditor.onDidChangeModelContent(htmlTracker.trackEdit);
      htmlEditor.onDidFocusEditorText(htmlTracker.trackFocus);
      htmlEditor.onDidChangeCursorSelection(htmlTracker.trackSelection);
    }

    if (goEditor) {
      const goTracker = this.createDebouncedTracker(goEditor, 'go');

      // Setup event listeners
      goEditor.onDidChangeModelContent(goTracker.trackEdit);
      goEditor.onDidFocusEditorText(goTracker.trackFocus);
      goEditor.onDidChangeCursorSelection(goTracker.trackSelection);
    }
  }
};

// Track app initialization when the page loads
document.addEventListener('DOMContentLoaded', function () {
  if (window.Analytics) {
    // window.Analytics.trackAppInit();
  }
}); 