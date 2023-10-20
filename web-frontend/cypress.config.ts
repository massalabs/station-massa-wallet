import { defineConfig } from 'cypress';
import getCompareSnapshotsPlugin from 'cypress-image-diff-js/dist/plugin';

export default defineConfig({
  component: {
    // macbook-15 default
    viewportWidth: 1440,
    viewportHeight: 900,
    devServer: {
      framework: 'react',
      bundler: 'vite',
    },
    setupNodeEvents(on, config) {
      const compareSnapshotsPlugin = getCompareSnapshotsPlugin(on, config);

      // Add the imageSnapshot configuration options to the plugin
      compareSnapshotsPlugin.options = {
        threshold: 0.05,
        thresholdType: 'pixel',
        includeAA: true,
        diffColor: [255, 0, 0],
        diffColorAlt: [0, 255, 0],
        capture: 'fullPage',
        antialiasingTolerance: 0,
        createDiffImage: true,
        diffPath: 'cypress/diffs',
        update: false,
        debug: false,
        errorColor: [255, 0, 0],
        errorType: 'movement',
        waitBeforeScreenshot: 0,
        waitAfterScreenshot: 0,
        tabbableOptions: {
          include: ['button', 'input', 'select', 'textarea', 'a[href]'],
        },
      };

      on('task', {
        deleteScreenshot() {
          // your task code here
          return null;
        },
      });

      return compareSnapshotsPlugin;
    },
  },
  env: {
    browserPermissions: {
      clipboard: 'allow',
    },
  },
  e2e: {
    // macbook-15 default
    viewportWidth: 1440,
    viewportHeight: 900,
    baseUrl: 'http://localhost:8080',
    setupNodeEvents(on, config) {
      const compareSnapshotsPlugin = getCompareSnapshotsPlugin(on, config);

      // Add the imageSnapshot configuration options to the plugin
      compareSnapshotsPlugin.options = {
        threshold: 0.05,
        thresholdType: 'pixel',
        includeAA: true,
        diffColor: [255, 0, 0],
        diffColorAlt: [0, 255, 0],
        capture: 'fullPage',
        antialiasingTolerance: 0,
        createDiffImage: true,
        diffPath: 'cypress/diffs',
        update: false,
        debug: false,
        errorColor: [255, 0, 0],
        errorType: 'movement',
        waitBeforeScreenshot: 0,
        waitAfterScreenshot: 0,
        tabbableOptions: {
          include: ['button', 'input', 'select', 'textarea', 'a[href]'],
        },
      };

      on('task', {
        deleteScreenshot() {
          // your task code here
          return null;
        },
      });

      return compareSnapshotsPlugin;
    },
  },
});
