// vite.config.ts
import { defineConfig, loadEnv } from "file:///D:/goCodes/devops-asset/web/node_modules/.pnpm/vite@5.4.19_@types+node@22.10.0_sass@1.81.0_terser@5.36.0/node_modules/vite/dist/node/index.js";
import vue from "file:///D:/goCodes/devops-asset/web/node_modules/.pnpm/@vitejs+plugin-vue@5.2.1_vite@5.4.19_@types+node@22.10.0_sass@1.81.0_terser@5.36.0__vue@3.5.15_typescript@5.6.3_/node_modules/@vitejs/plugin-vue/dist/index.mjs";
import path from "path";
import viteCompression from "file:///D:/goCodes/devops-asset/web/node_modules/.pnpm/vite-plugin-compression@0.5.1_vite@5.4.19_@types+node@22.10.0_sass@1.81.0_terser@5.36.0_/node_modules/vite-plugin-compression/dist/index.mjs";
import Components from "file:///D:/goCodes/devops-asset/web/node_modules/.pnpm/unplugin-vue-components@0.26.0_@babel+parser@7.27.5_@nuxt+kit@3.17.4_rollup@4.34.8_vue@3.5.15_typescript@5.6.3_/node_modules/unplugin-vue-components/dist/vite.js";
import AutoImport from "file:///D:/goCodes/devops-asset/web/node_modules/.pnpm/unplugin-auto-import@0.18.5_@nuxt+kit@3.17.4_@vueuse+core@11.3.0_vue@3.5.15_typescript@5.6.3___rollup@4.34.8/node_modules/unplugin-auto-import/dist/vite.js";
import { ElementPlusResolver } from "file:///D:/goCodes/devops-asset/web/node_modules/.pnpm/unplugin-vue-components@0.26.0_@babel+parser@7.27.5_@nuxt+kit@3.17.4_rollup@4.34.8_vue@3.5.15_typescript@5.6.3_/node_modules/unplugin-vue-components/dist/resolvers.js";
import { fileURLToPath } from "url";
import vueDevTools from "file:///D:/goCodes/devops-asset/web/node_modules/.pnpm/vite-plugin-vue-devtools@7.7.6_@nuxt+kit@3.17.4_rollup@4.34.8_vite@5.4.19_@types+node@22.10.0_eotppr3nkkktieujjhcnq4nsl4/node_modules/vite-plugin-vue-devtools/dist/vite.mjs";
var __vite_injected_original_dirname = "D:\\goCodes\\devops-asset\\web";
var __vite_injected_original_import_meta_url = "file:///D:/goCodes/devops-asset/web/vite.config.ts";
var vite_config_default = ({ mode }) => {
  const root = process.cwd();
  const env = loadEnv(mode, root);
  const { VITE_VERSION, VITE_PORT, VITE_BASE_URL, VITE_API_URL } = env;
  console.log(`\u{1F680} API_URL = ${VITE_API_URL}`);
  console.log(`\u{1F680} VERSION = ${VITE_VERSION}`);
  return defineConfig({
    define: {
      __APP_VERSION__: JSON.stringify(VITE_VERSION)
    },
    base: VITE_BASE_URL,
    server: {
      port: parseInt(VITE_PORT),
      proxy: {
        "/api": {
          target: VITE_API_URL,
          changeOrigin: true,
          rewrite: (path2) => path2.replace(/^\/api/, "/api")
        }
      },
      host: true
    },
    // 路径别名
    resolve: {
      alias: {
        "@": fileURLToPath(new URL("./src", __vite_injected_original_import_meta_url)),
        "@views": resolvePath("src/views"),
        "@imgs": resolvePath("src/assets/img"),
        "@icons": resolvePath("src/assets/icons"),
        "@utils": resolvePath("src/utils"),
        "@stores": resolvePath("src/store"),
        "@plugins": resolvePath("src/plugins"),
        "@styles": resolvePath("src/assets/styles")
      }
    },
    build: {
      target: "es2015",
      outDir: "dist",
      chunkSizeWarningLimit: 2e3,
      minify: "terser",
      terserOptions: {
        compress: {
          drop_console: true,
          // 生产环境去除 console
          drop_debugger: true
          // 生产环境去除 debugger
        }
      },
      rollupOptions: {
        output: {
          manualChunks: {
            vendor: ["vue", "vue-router", "pinia", "element-plus"]
          }
        }
      },
      dynamicImportVarsOptions: {
        warnOnError: true,
        exclude: [],
        include: ["src/views/**/*.vue"]
      }
    },
    plugins: [
      vue(),
      // 自动导入 components 下面的组件，无需 import 引入
      Components({
        deep: true,
        extensions: ["vue"],
        dirs: ["src/components"],
        // 自动导入的组件目录
        resolvers: [ElementPlusResolver()],
        dts: "src/types/components.d.ts"
        // 指定类型声明文件的路径
      }),
      AutoImport({
        imports: ["vue", "vue-router", "@vueuse/core", "pinia"],
        resolvers: [ElementPlusResolver()],
        dts: "src/types/auto-imports.d.ts",
        eslintrc: {
          // 这里先设置成true然后pnpm dev 运行之后会生成 .auto-import.json 文件之后，在改为false
          enabled: true,
          filepath: "./.auto-import.json",
          globalsPropValue: true
        }
      }),
      // 打包分析
      // visualizer({
      //   open: true,
      //   gzipSize: true,
      //   brotliSize: true,
      //   filename: 'dist/stats.html' // 分析图生成的文件名及路径
      // }),
      // 压缩
      viteCompression({
        verbose: true,
        // 是否在控制台输出压缩结果
        disable: false,
        // 是否禁用
        algorithm: "gzip",
        // 压缩算法,可选 [ 'gzip' , 'brotliCompress' ,'deflate' , 'deflateRaw']
        ext: ".gz",
        // 压缩后的文件名后缀
        threshold: 10240,
        // 只有大小大于该值的资源会被处理 10240B = 10KB
        deleteOriginFile: false
        // 压缩后是否删除原文件
      }),
      // 图片压缩
      // viteImagemin({
      //   verbose: true, // 是否在控制台输出压缩结果
      //   // 图片压缩配置
      //   // GIF 图片压缩配置
      //   gifsicle: {
      //     optimizationLevel: 4, // 优化级别 1-7，7为最高级别压缩
      //     interlaced: false // 是否隔行扫描
      //   },
      //   // PNG 图片压缩配置
      //   optipng: {
      //     optimizationLevel: 4 // 优化级别 0-7，7为最高级别压缩
      //   },
      //   // JPEG 图片压缩配置
      //   mozjpeg: {
      //     quality: 60 // 压缩质量 0-100，值越小压缩率越高
      //   },
      //   // PNG 图片压缩配置(另一个压缩器)
      //   pngquant: {
      //     quality: [0.8, 0.9], // 压缩质量范围 0-1
      //     speed: 4 // 压缩速度 1-11，值越大压缩速度越快，但质量可能会下降
      //   },
      //   // SVG 图片压缩配置
      //   svgo: {
      //     plugins: [
      //       {
      //         name: 'removeViewBox' // 移除 viewBox 属性
      //       },
      //       {
      //         name: 'removeEmptyAttrs', // 移除空属性
      //         active: false // 是否启用此插件
      //       }
      //     ]
      //   }
      // })
      vueDevTools()
    ],
    // 预加载项目必需的组件
    optimizeDeps: {
      include: [
        "vue",
        "vue-router",
        "pinia",
        "axios",
        "@vueuse/core",
        "echarts",
        "@wangeditor/editor",
        "@wangeditor/editor-for-vue",
        "vue-i18n",
        "element-plus/es/components/form/style/css",
        "element-plus/es/components/form-item/style/css",
        "element-plus/es/components/button/style/css",
        "element-plus/es/components/input/style/css",
        "element-plus/es/components/input-number/style/css",
        "element-plus/es/components/switch/style/css",
        "element-plus/es/components/upload/style/css",
        "element-plus/es/components/menu/style/css",
        "element-plus/es/components/col/style/css",
        "element-plus/es/components/icon/style/css",
        "element-plus/es/components/row/style/css",
        "element-plus/es/components/tag/style/css",
        "element-plus/es/components/dialog/style/css",
        "element-plus/es/components/loading/style/css",
        "element-plus/es/components/radio/style/css",
        "element-plus/es/components/radio-group/style/css",
        "element-plus/es/components/popover/style/css",
        "element-plus/es/components/scrollbar/style/css",
        "element-plus/es/components/tooltip/style/css",
        "element-plus/es/components/dropdown/style/css",
        "element-plus/es/components/dropdown-menu/style/css",
        "element-plus/es/components/dropdown-item/style/css",
        "element-plus/es/components/sub-menu/style/css",
        "element-plus/es/components/menu-item/style/css",
        "element-plus/es/components/divider/style/css",
        "element-plus/es/components/card/style/css",
        "element-plus/es/components/link/style/css",
        "element-plus/es/components/breadcrumb/style/css",
        "element-plus/es/components/breadcrumb-item/style/css",
        "element-plus/es/components/table/style/css",
        "element-plus/es/components/tree-select/style/css",
        "element-plus/es/components/table-column/style/css",
        "element-plus/es/components/select/style/css",
        "element-plus/es/components/option/style/css",
        "element-plus/es/components/pagination/style/css",
        "element-plus/es/components/tree/style/css",
        "element-plus/es/components/alert/style/css",
        "element-plus/es/components/radio-button/style/css",
        "element-plus/es/components/checkbox-group/style/css",
        "element-plus/es/components/checkbox/style/css",
        "element-plus/es/components/tabs/style/css",
        "element-plus/es/components/tab-pane/style/css",
        "element-plus/es/components/rate/style/css",
        "element-plus/es/components/date-picker/style/css",
        "element-plus/es/components/notification/style/css",
        "element-plus/es/components/image/style/css",
        "element-plus/es/components/statistic/style/css",
        "element-plus/es/components/watermark/style/css",
        "element-plus/es/components/config-provider/style/css",
        "element-plus/es/components/text/style/css",
        "element-plus/es/components/drawer/style/css",
        "element-plus/es/components/color-picker/style/css",
        "element-plus/es/components/backtop/style/css",
        "element-plus/es/components/message-box/style/css",
        "element-plus/es/components/skeleton/style/css",
        "element-plus/es/components/skeleton/style/css",
        "element-plus/es/components/skeleton-item/style/css",
        "element-plus/es/components/badge/style/css",
        "element-plus/es/components/steps/style/css",
        "element-plus/es/components/step/style/css",
        "element-plus/es/components/avatar/style/css",
        "element-plus/es/components/descriptions/style/css",
        "element-plus/es/components/descriptions-item/style/css",
        "element-plus/es/components/checkbox-group/style/css",
        "element-plus/es/components/progress/style/css",
        "element-plus/es/components/image-viewer/style/css",
        "element-plus/es/components/empty/style/css",
        "element-plus/es/components/segmented/style/css",
        "element-plus/es/components/calendar/style/css",
        "element-plus/es/components/message/style/css",
        "xlsx",
        "file-saver",
        "element-plus/es/components/timeline/style/css",
        "element-plus/es/components/timeline-item/style/css",
        "vue-img-cutter"
      ]
    },
    css: {
      preprocessorOptions: {
        // sass variable and mixin
        scss: {
          api: "modern-compiler",
          additionalData: `
            @use "@styles/variables.scss" as *; @use "@styles/mixin.scss" as *;
          `
        }
      },
      postcss: {
        plugins: [
          {
            postcssPlugin: "internal:charset-removal",
            AtRule: {
              charset: (atRule) => {
                if (atRule.name === "charset") {
                  atRule.remove();
                }
              }
            }
          }
        ]
      }
    }
  });
};
function resolvePath(paths) {
  return path.resolve(__vite_injected_original_dirname, paths);
}
export {
  vite_config_default as default
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcudHMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbImNvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9kaXJuYW1lID0gXCJEOlxcXFxnb0NvZGVzXFxcXGRldm9wcy1hc3NldFxcXFx3ZWJcIjtjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfZmlsZW5hbWUgPSBcIkQ6XFxcXGdvQ29kZXNcXFxcZGV2b3BzLWFzc2V0XFxcXHdlYlxcXFx2aXRlLmNvbmZpZy50c1wiO2NvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9pbXBvcnRfbWV0YV91cmwgPSBcImZpbGU6Ly8vRDovZ29Db2Rlcy9kZXZvcHMtYXNzZXQvd2ViL3ZpdGUuY29uZmlnLnRzXCI7aW1wb3J0IHsgZGVmaW5lQ29uZmlnLCBsb2FkRW52IH0gZnJvbSAndml0ZSdcclxuaW1wb3J0IHZ1ZSBmcm9tICdAdml0ZWpzL3BsdWdpbi12dWUnXHJcbmltcG9ydCBwYXRoIGZyb20gJ3BhdGgnXHJcbmltcG9ydCB2aXRlQ29tcHJlc3Npb24gZnJvbSAndml0ZS1wbHVnaW4tY29tcHJlc3Npb24nXHJcbmltcG9ydCBDb21wb25lbnRzIGZyb20gJ3VucGx1Z2luLXZ1ZS1jb21wb25lbnRzL3ZpdGUnXHJcbmltcG9ydCBBdXRvSW1wb3J0IGZyb20gJ3VucGx1Z2luLWF1dG8taW1wb3J0L3ZpdGUnXHJcbmltcG9ydCB7IEVsZW1lbnRQbHVzUmVzb2x2ZXIgfSBmcm9tICd1bnBsdWdpbi12dWUtY29tcG9uZW50cy9yZXNvbHZlcnMnXHJcbmltcG9ydCB7IGZpbGVVUkxUb1BhdGggfSBmcm9tICd1cmwnXHJcbi8vIGltcG9ydCB2aXRlSW1hZ2VtaW4gZnJvbSAndml0ZS1wbHVnaW4taW1hZ2VtaW4nXHJcbi8vIGltcG9ydCB7IHZpc3VhbGl6ZXIgfSBmcm9tICdyb2xsdXAtcGx1Z2luLXZpc3VhbGl6ZXInXHJcblxyXG4vLyBodHRwczovL2RldnRvb2xzLnZ1ZWpzLm9yZy9nZXR0aW5nLXN0YXJ0ZWQvaW50cm9kdWN0aW9uXHJcbmltcG9ydCB2dWVEZXZUb29scyBmcm9tICd2aXRlLXBsdWdpbi12dWUtZGV2dG9vbHMnXHJcblxyXG5leHBvcnQgZGVmYXVsdCAoeyBtb2RlIH06IHsgbW9kZTogc3RyaW5nIH0pID0+IHtcclxuICBjb25zdCByb290ID0gcHJvY2Vzcy5jd2QoKVxyXG4gIGNvbnN0IGVudiA9IGxvYWRFbnYobW9kZSwgcm9vdClcclxuICBjb25zdCB7IFZJVEVfVkVSU0lPTiwgVklURV9QT1JULCBWSVRFX0JBU0VfVVJMLCBWSVRFX0FQSV9VUkwgfSA9IGVudlxyXG5cclxuICBjb25zb2xlLmxvZyhgXHVEODNEXHVERTgwIEFQSV9VUkwgPSAke1ZJVEVfQVBJX1VSTH1gKVxyXG4gIGNvbnNvbGUubG9nKGBcdUQ4M0RcdURFODAgVkVSU0lPTiA9ICR7VklURV9WRVJTSU9OfWApXHJcblxyXG4gIHJldHVybiBkZWZpbmVDb25maWcoe1xyXG4gICAgZGVmaW5lOiB7XHJcbiAgICAgIF9fQVBQX1ZFUlNJT05fXzogSlNPTi5zdHJpbmdpZnkoVklURV9WRVJTSU9OKVxyXG4gICAgfSxcclxuICAgIGJhc2U6IFZJVEVfQkFTRV9VUkwsXHJcbiAgICBzZXJ2ZXI6IHtcclxuICAgICAgcG9ydDogcGFyc2VJbnQoVklURV9QT1JUKSxcclxuICAgICAgcHJveHk6IHtcclxuICAgICAgICAnL2FwaSc6IHtcclxuICAgICAgICAgIHRhcmdldDogVklURV9BUElfVVJMLFxyXG4gICAgICAgICAgY2hhbmdlT3JpZ2luOiB0cnVlLFxyXG4gICAgICAgICAgcmV3cml0ZTogKHBhdGgpID0+IHBhdGgucmVwbGFjZSgvXlxcL2FwaS8sICcvYXBpJylcclxuICAgICAgICB9XHJcbiAgICAgIH0sXHJcbiAgICAgIGhvc3Q6IHRydWVcclxuICAgIH0sXHJcbiAgICAvLyBcdThERUZcdTVGODRcdTUyMkJcdTU0MERcclxuICAgIHJlc29sdmU6IHtcclxuICAgICAgYWxpYXM6IHtcclxuICAgICAgICAnQCc6IGZpbGVVUkxUb1BhdGgobmV3IFVSTCgnLi9zcmMnLCBpbXBvcnQubWV0YS51cmwpKSxcclxuICAgICAgICAnQHZpZXdzJzogcmVzb2x2ZVBhdGgoJ3NyYy92aWV3cycpLFxyXG4gICAgICAgICdAaW1ncyc6IHJlc29sdmVQYXRoKCdzcmMvYXNzZXRzL2ltZycpLFxyXG4gICAgICAgICdAaWNvbnMnOiByZXNvbHZlUGF0aCgnc3JjL2Fzc2V0cy9pY29ucycpLFxyXG4gICAgICAgICdAdXRpbHMnOiByZXNvbHZlUGF0aCgnc3JjL3V0aWxzJyksXHJcbiAgICAgICAgJ0BzdG9yZXMnOiByZXNvbHZlUGF0aCgnc3JjL3N0b3JlJyksXHJcbiAgICAgICAgJ0BwbHVnaW5zJzogcmVzb2x2ZVBhdGgoJ3NyYy9wbHVnaW5zJyksXHJcbiAgICAgICAgJ0BzdHlsZXMnOiByZXNvbHZlUGF0aCgnc3JjL2Fzc2V0cy9zdHlsZXMnKVxyXG4gICAgICB9XHJcbiAgICB9LFxyXG4gICAgYnVpbGQ6IHtcclxuICAgICAgdGFyZ2V0OiAnZXMyMDE1JyxcclxuICAgICAgb3V0RGlyOiAnZGlzdCcsXHJcbiAgICAgIGNodW5rU2l6ZVdhcm5pbmdMaW1pdDogMjAwMCxcclxuICAgICAgbWluaWZ5OiAndGVyc2VyJyxcclxuICAgICAgdGVyc2VyT3B0aW9uczoge1xyXG4gICAgICAgIGNvbXByZXNzOiB7XHJcbiAgICAgICAgICBkcm9wX2NvbnNvbGU6IHRydWUsIC8vIFx1NzUxRlx1NEVBN1x1NzNBRlx1NTg4M1x1NTNCQlx1OTY2NCBjb25zb2xlXHJcbiAgICAgICAgICBkcm9wX2RlYnVnZ2VyOiB0cnVlIC8vIFx1NzUxRlx1NEVBN1x1NzNBRlx1NTg4M1x1NTNCQlx1OTY2NCBkZWJ1Z2dlclxyXG4gICAgICAgIH1cclxuICAgICAgfSxcclxuICAgICAgcm9sbHVwT3B0aW9uczoge1xyXG4gICAgICAgIG91dHB1dDoge1xyXG4gICAgICAgICAgbWFudWFsQ2h1bmtzOiB7XHJcbiAgICAgICAgICAgIHZlbmRvcjogWyd2dWUnLCAndnVlLXJvdXRlcicsICdwaW5pYScsICdlbGVtZW50LXBsdXMnXVxyXG4gICAgICAgICAgfVxyXG4gICAgICAgIH1cclxuICAgICAgfSxcclxuICAgICAgZHluYW1pY0ltcG9ydFZhcnNPcHRpb25zOiB7XHJcbiAgICAgICAgd2Fybk9uRXJyb3I6IHRydWUsXHJcbiAgICAgICAgZXhjbHVkZTogW10sXHJcbiAgICAgICAgaW5jbHVkZTogWydzcmMvdmlld3MvKiovKi52dWUnXVxyXG4gICAgICB9XHJcbiAgICB9LFxyXG4gICAgcGx1Z2luczogW1xyXG4gICAgICB2dWUoKSxcclxuICAgICAgLy8gXHU4MUVBXHU1MkE4XHU1QkZDXHU1MTY1IGNvbXBvbmVudHMgXHU0RTBCXHU5NzYyXHU3Njg0XHU3RUM0XHU0RUY2XHVGRjBDXHU2NUUwXHU5NzAwIGltcG9ydCBcdTVGMTVcdTUxNjVcclxuICAgICAgQ29tcG9uZW50cyh7XHJcbiAgICAgICAgZGVlcDogdHJ1ZSxcclxuICAgICAgICBleHRlbnNpb25zOiBbJ3Z1ZSddLFxyXG4gICAgICAgIGRpcnM6IFsnc3JjL2NvbXBvbmVudHMnXSwgLy8gXHU4MUVBXHU1MkE4XHU1QkZDXHU1MTY1XHU3Njg0XHU3RUM0XHU0RUY2XHU3NkVFXHU1RjU1XHJcbiAgICAgICAgcmVzb2x2ZXJzOiBbRWxlbWVudFBsdXNSZXNvbHZlcigpXSxcclxuICAgICAgICBkdHM6ICdzcmMvdHlwZXMvY29tcG9uZW50cy5kLnRzJyAvLyBcdTYzMDdcdTVCOUFcdTdDN0JcdTU3OEJcdTU4RjBcdTY2MEVcdTY1ODdcdTRFRjZcdTc2ODRcdThERUZcdTVGODRcclxuICAgICAgfSksXHJcbiAgICAgIEF1dG9JbXBvcnQoe1xyXG4gICAgICAgIGltcG9ydHM6IFsndnVlJywgJ3Z1ZS1yb3V0ZXInLCAnQHZ1ZXVzZS9jb3JlJywgJ3BpbmlhJ10sXHJcbiAgICAgICAgcmVzb2x2ZXJzOiBbRWxlbWVudFBsdXNSZXNvbHZlcigpXSxcclxuICAgICAgICBkdHM6ICdzcmMvdHlwZXMvYXV0by1pbXBvcnRzLmQudHMnLFxyXG4gICAgICAgIGVzbGludHJjOiB7XHJcbiAgICAgICAgICAvLyBcdThGRDlcdTkxQ0NcdTUxNDhcdThCQkVcdTdGNkVcdTYyMTB0cnVlXHU3MTM2XHU1NDBFcG5wbSBkZXYgXHU4RkQwXHU4ODRDXHU0RTRCXHU1NDBFXHU0RjFBXHU3NTFGXHU2MjEwIC5hdXRvLWltcG9ydC5qc29uIFx1NjU4N1x1NEVGNlx1NEU0Qlx1NTQwRVx1RkYwQ1x1NTcyOFx1NjUzOVx1NEUzQWZhbHNlXHJcbiAgICAgICAgICBlbmFibGVkOiB0cnVlLFxyXG4gICAgICAgICAgZmlsZXBhdGg6ICcuLy5hdXRvLWltcG9ydC5qc29uJyxcclxuICAgICAgICAgIGdsb2JhbHNQcm9wVmFsdWU6IHRydWVcclxuICAgICAgICB9XHJcbiAgICAgIH0pLFxyXG4gICAgICAvLyBcdTYyNTNcdTUzMDVcdTUyMDZcdTY3OTBcclxuICAgICAgLy8gdmlzdWFsaXplcih7XHJcbiAgICAgIC8vICAgb3BlbjogdHJ1ZSxcclxuICAgICAgLy8gICBnemlwU2l6ZTogdHJ1ZSxcclxuICAgICAgLy8gICBicm90bGlTaXplOiB0cnVlLFxyXG4gICAgICAvLyAgIGZpbGVuYW1lOiAnZGlzdC9zdGF0cy5odG1sJyAvLyBcdTUyMDZcdTY3OTBcdTU2RkVcdTc1MUZcdTYyMTBcdTc2ODRcdTY1ODdcdTRFRjZcdTU0MERcdTUzQ0FcdThERUZcdTVGODRcclxuICAgICAgLy8gfSksXHJcbiAgICAgIC8vIFx1NTM4Qlx1N0YyOVxyXG4gICAgICB2aXRlQ29tcHJlc3Npb24oe1xyXG4gICAgICAgIHZlcmJvc2U6IHRydWUsIC8vIFx1NjYyRlx1NTQyNlx1NTcyOFx1NjNBN1x1NTIzNlx1NTNGMFx1OEY5M1x1NTFGQVx1NTM4Qlx1N0YyOVx1N0VEM1x1Njc5Q1xyXG4gICAgICAgIGRpc2FibGU6IGZhbHNlLCAvLyBcdTY2MkZcdTU0MjZcdTc5ODFcdTc1MjhcclxuICAgICAgICBhbGdvcml0aG06ICdnemlwJywgLy8gXHU1MzhCXHU3RjI5XHU3Qjk3XHU2Q0Q1LFx1NTNFRlx1OTAwOSBbICdnemlwJyAsICdicm90bGlDb21wcmVzcycgLCdkZWZsYXRlJyAsICdkZWZsYXRlUmF3J11cclxuICAgICAgICBleHQ6ICcuZ3onLCAvLyBcdTUzOEJcdTdGMjlcdTU0MEVcdTc2ODRcdTY1ODdcdTRFRjZcdTU0MERcdTU0MEVcdTdGMDBcclxuICAgICAgICB0aHJlc2hvbGQ6IDEwMjQwLCAvLyBcdTUzRUFcdTY3MDlcdTU5MjdcdTVDMEZcdTU5MjdcdTRFOEVcdThCRTVcdTUwM0NcdTc2ODRcdThENDRcdTZFOTBcdTRGMUFcdTg4QUJcdTU5MDRcdTc0MDYgMTAyNDBCID0gMTBLQlxyXG4gICAgICAgIGRlbGV0ZU9yaWdpbkZpbGU6IGZhbHNlIC8vIFx1NTM4Qlx1N0YyOVx1NTQwRVx1NjYyRlx1NTQyNlx1NTIyMFx1OTY2NFx1NTM5Rlx1NjU4N1x1NEVGNlxyXG4gICAgICB9KSxcclxuICAgICAgLy8gXHU1NkZFXHU3MjQ3XHU1MzhCXHU3RjI5XHJcbiAgICAgIC8vIHZpdGVJbWFnZW1pbih7XHJcbiAgICAgIC8vICAgdmVyYm9zZTogdHJ1ZSwgLy8gXHU2NjJGXHU1NDI2XHU1NzI4XHU2M0E3XHU1MjM2XHU1M0YwXHU4RjkzXHU1MUZBXHU1MzhCXHU3RjI5XHU3RUQzXHU2NzlDXHJcbiAgICAgIC8vICAgLy8gXHU1NkZFXHU3MjQ3XHU1MzhCXHU3RjI5XHU5MTREXHU3RjZFXHJcbiAgICAgIC8vICAgLy8gR0lGIFx1NTZGRVx1NzI0N1x1NTM4Qlx1N0YyOVx1OTE0RFx1N0Y2RVxyXG4gICAgICAvLyAgIGdpZnNpY2xlOiB7XHJcbiAgICAgIC8vICAgICBvcHRpbWl6YXRpb25MZXZlbDogNCwgLy8gXHU0RjE4XHU1MzE2XHU3RUE3XHU1MjJCIDEtN1x1RkYwQzdcdTRFM0FcdTY3MDBcdTlBRDhcdTdFQTdcdTUyMkJcdTUzOEJcdTdGMjlcclxuICAgICAgLy8gICAgIGludGVybGFjZWQ6IGZhbHNlIC8vIFx1NjYyRlx1NTQyNlx1OTY5NFx1ODg0Q1x1NjI2Qlx1NjNDRlxyXG4gICAgICAvLyAgIH0sXHJcbiAgICAgIC8vICAgLy8gUE5HIFx1NTZGRVx1NzI0N1x1NTM4Qlx1N0YyOVx1OTE0RFx1N0Y2RVxyXG4gICAgICAvLyAgIG9wdGlwbmc6IHtcclxuICAgICAgLy8gICAgIG9wdGltaXphdGlvbkxldmVsOiA0IC8vIFx1NEYxOFx1NTMxNlx1N0VBN1x1NTIyQiAwLTdcdUZGMEM3XHU0RTNBXHU2NzAwXHU5QUQ4XHU3RUE3XHU1MjJCXHU1MzhCXHU3RjI5XHJcbiAgICAgIC8vICAgfSxcclxuICAgICAgLy8gICAvLyBKUEVHIFx1NTZGRVx1NzI0N1x1NTM4Qlx1N0YyOVx1OTE0RFx1N0Y2RVxyXG4gICAgICAvLyAgIG1vempwZWc6IHtcclxuICAgICAgLy8gICAgIHF1YWxpdHk6IDYwIC8vIFx1NTM4Qlx1N0YyOVx1OEQyOFx1OTFDRiAwLTEwMFx1RkYwQ1x1NTAzQ1x1OEQ4QVx1NUMwRlx1NTM4Qlx1N0YyOVx1NzM4N1x1OEQ4QVx1OUFEOFxyXG4gICAgICAvLyAgIH0sXHJcbiAgICAgIC8vICAgLy8gUE5HIFx1NTZGRVx1NzI0N1x1NTM4Qlx1N0YyOVx1OTE0RFx1N0Y2RShcdTUzRTZcdTRFMDBcdTRFMkFcdTUzOEJcdTdGMjlcdTU2NjgpXHJcbiAgICAgIC8vICAgcG5ncXVhbnQ6IHtcclxuICAgICAgLy8gICAgIHF1YWxpdHk6IFswLjgsIDAuOV0sIC8vIFx1NTM4Qlx1N0YyOVx1OEQyOFx1OTFDRlx1ODMwM1x1NTZGNCAwLTFcclxuICAgICAgLy8gICAgIHNwZWVkOiA0IC8vIFx1NTM4Qlx1N0YyOVx1OTAxRlx1NUVBNiAxLTExXHVGRjBDXHU1MDNDXHU4RDhBXHU1OTI3XHU1MzhCXHU3RjI5XHU5MDFGXHU1RUE2XHU4RDhBXHU1RkVCXHVGRjBDXHU0RjQ2XHU4RDI4XHU5MUNGXHU1M0VGXHU4MEZEXHU0RjFBXHU0RTBCXHU5NjREXHJcbiAgICAgIC8vICAgfSxcclxuICAgICAgLy8gICAvLyBTVkcgXHU1NkZFXHU3MjQ3XHU1MzhCXHU3RjI5XHU5MTREXHU3RjZFXHJcbiAgICAgIC8vICAgc3Znbzoge1xyXG4gICAgICAvLyAgICAgcGx1Z2luczogW1xyXG4gICAgICAvLyAgICAgICB7XHJcbiAgICAgIC8vICAgICAgICAgbmFtZTogJ3JlbW92ZVZpZXdCb3gnIC8vIFx1NzlGQlx1OTY2NCB2aWV3Qm94IFx1NUM1RVx1NjAyN1xyXG4gICAgICAvLyAgICAgICB9LFxyXG4gICAgICAvLyAgICAgICB7XHJcbiAgICAgIC8vICAgICAgICAgbmFtZTogJ3JlbW92ZUVtcHR5QXR0cnMnLCAvLyBcdTc5RkJcdTk2NjRcdTdBN0FcdTVDNUVcdTYwMjdcclxuICAgICAgLy8gICAgICAgICBhY3RpdmU6IGZhbHNlIC8vIFx1NjYyRlx1NTQyNlx1NTQyRlx1NzUyOFx1NkI2NFx1NjNEMlx1NEVGNlxyXG4gICAgICAvLyAgICAgICB9XHJcbiAgICAgIC8vICAgICBdXHJcbiAgICAgIC8vICAgfVxyXG4gICAgICAvLyB9KVxyXG4gICAgICB2dWVEZXZUb29scygpXHJcbiAgICBdLFxyXG4gICAgLy8gXHU5ODg0XHU1MkEwXHU4RjdEXHU5ODc5XHU3NkVFXHU1RkM1XHU5NzAwXHU3Njg0XHU3RUM0XHU0RUY2XHJcbiAgICBvcHRpbWl6ZURlcHM6IHtcclxuICAgICAgaW5jbHVkZTogW1xyXG4gICAgICAgICd2dWUnLFxyXG4gICAgICAgICd2dWUtcm91dGVyJyxcclxuICAgICAgICAncGluaWEnLFxyXG4gICAgICAgICdheGlvcycsXHJcbiAgICAgICAgJ0B2dWV1c2UvY29yZScsXHJcbiAgICAgICAgJ2VjaGFydHMnLFxyXG4gICAgICAgICdAd2FuZ2VkaXRvci9lZGl0b3InLFxyXG4gICAgICAgICdAd2FuZ2VkaXRvci9lZGl0b3ItZm9yLXZ1ZScsXHJcbiAgICAgICAgJ3Z1ZS1pMThuJyxcclxuICAgICAgICAnZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvZm9ybS9zdHlsZS9jc3MnLFxyXG4gICAgICAgICdlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9mb3JtLWl0ZW0vc3R5bGUvY3NzJyxcclxuICAgICAgICAnZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvYnV0dG9uL3N0eWxlL2NzcycsXHJcbiAgICAgICAgJ2VsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL2lucHV0L3N0eWxlL2NzcycsXHJcbiAgICAgICAgJ2VsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL2lucHV0LW51bWJlci9zdHlsZS9jc3MnLFxyXG4gICAgICAgICdlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9zd2l0Y2gvc3R5bGUvY3NzJyxcclxuICAgICAgICAnZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvdXBsb2FkL3N0eWxlL2NzcycsXHJcbiAgICAgICAgJ2VsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL21lbnUvc3R5bGUvY3NzJyxcclxuICAgICAgICAnZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvY29sL3N0eWxlL2NzcycsXHJcbiAgICAgICAgJ2VsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL2ljb24vc3R5bGUvY3NzJyxcclxuICAgICAgICAnZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvcm93L3N0eWxlL2NzcycsXHJcbiAgICAgICAgJ2VsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL3RhZy9zdHlsZS9jc3MnLFxyXG4gICAgICAgICdlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9kaWFsb2cvc3R5bGUvY3NzJyxcclxuICAgICAgICAnZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvbG9hZGluZy9zdHlsZS9jc3MnLFxyXG4gICAgICAgICdlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9yYWRpby9zdHlsZS9jc3MnLFxyXG4gICAgICAgICdlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9yYWRpby1ncm91cC9zdHlsZS9jc3MnLFxyXG4gICAgICAgICdlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9wb3BvdmVyL3N0eWxlL2NzcycsXHJcbiAgICAgICAgJ2VsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL3Njcm9sbGJhci9zdHlsZS9jc3MnLFxyXG4gICAgICAgICdlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy90b29sdGlwL3N0eWxlL2NzcycsXHJcbiAgICAgICAgJ2VsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL2Ryb3Bkb3duL3N0eWxlL2NzcycsXHJcbiAgICAgICAgJ2VsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL2Ryb3Bkb3duLW1lbnUvc3R5bGUvY3NzJyxcclxuICAgICAgICAnZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvZHJvcGRvd24taXRlbS9zdHlsZS9jc3MnLFxyXG4gICAgICAgICdlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9zdWItbWVudS9zdHlsZS9jc3MnLFxyXG4gICAgICAgICdlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9tZW51LWl0ZW0vc3R5bGUvY3NzJyxcclxuICAgICAgICAnZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvZGl2aWRlci9zdHlsZS9jc3MnLFxyXG4gICAgICAgICdlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9jYXJkL3N0eWxlL2NzcycsXHJcbiAgICAgICAgJ2VsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL2xpbmsvc3R5bGUvY3NzJyxcclxuICAgICAgICAnZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvYnJlYWRjcnVtYi9zdHlsZS9jc3MnLFxyXG4gICAgICAgICdlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9icmVhZGNydW1iLWl0ZW0vc3R5bGUvY3NzJyxcclxuICAgICAgICAnZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvdGFibGUvc3R5bGUvY3NzJyxcclxuICAgICAgICAnZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvdHJlZS1zZWxlY3Qvc3R5bGUvY3NzJyxcclxuICAgICAgICAnZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvdGFibGUtY29sdW1uL3N0eWxlL2NzcycsXHJcbiAgICAgICAgJ2VsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL3NlbGVjdC9zdHlsZS9jc3MnLFxyXG4gICAgICAgICdlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9vcHRpb24vc3R5bGUvY3NzJyxcclxuICAgICAgICAnZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvcGFnaW5hdGlvbi9zdHlsZS9jc3MnLFxyXG4gICAgICAgICdlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy90cmVlL3N0eWxlL2NzcycsXHJcbiAgICAgICAgJ2VsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL2FsZXJ0L3N0eWxlL2NzcycsXHJcbiAgICAgICAgJ2VsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL3JhZGlvLWJ1dHRvbi9zdHlsZS9jc3MnLFxyXG4gICAgICAgICdlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9jaGVja2JveC1ncm91cC9zdHlsZS9jc3MnLFxyXG4gICAgICAgICdlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9jaGVja2JveC9zdHlsZS9jc3MnLFxyXG4gICAgICAgICdlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy90YWJzL3N0eWxlL2NzcycsXHJcbiAgICAgICAgJ2VsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL3RhYi1wYW5lL3N0eWxlL2NzcycsXHJcbiAgICAgICAgJ2VsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL3JhdGUvc3R5bGUvY3NzJyxcclxuICAgICAgICAnZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvZGF0ZS1waWNrZXIvc3R5bGUvY3NzJyxcclxuICAgICAgICAnZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvbm90aWZpY2F0aW9uL3N0eWxlL2NzcycsXHJcbiAgICAgICAgJ2VsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL2ltYWdlL3N0eWxlL2NzcycsXHJcbiAgICAgICAgJ2VsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL3N0YXRpc3RpYy9zdHlsZS9jc3MnLFxyXG4gICAgICAgICdlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy93YXRlcm1hcmsvc3R5bGUvY3NzJyxcclxuICAgICAgICAnZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvY29uZmlnLXByb3ZpZGVyL3N0eWxlL2NzcycsXHJcbiAgICAgICAgJ2VsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL3RleHQvc3R5bGUvY3NzJyxcclxuICAgICAgICAnZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvZHJhd2VyL3N0eWxlL2NzcycsXHJcbiAgICAgICAgJ2VsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL2NvbG9yLXBpY2tlci9zdHlsZS9jc3MnLFxyXG4gICAgICAgICdlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9iYWNrdG9wL3N0eWxlL2NzcycsXHJcbiAgICAgICAgJ2VsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL21lc3NhZ2UtYm94L3N0eWxlL2NzcycsXHJcbiAgICAgICAgJ2VsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL3NrZWxldG9uL3N0eWxlL2NzcycsXHJcbiAgICAgICAgJ2VsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL3NrZWxldG9uL3N0eWxlL2NzcycsXHJcbiAgICAgICAgJ2VsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL3NrZWxldG9uLWl0ZW0vc3R5bGUvY3NzJyxcclxuICAgICAgICAnZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvYmFkZ2Uvc3R5bGUvY3NzJyxcclxuICAgICAgICAnZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvc3RlcHMvc3R5bGUvY3NzJyxcclxuICAgICAgICAnZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvc3RlcC9zdHlsZS9jc3MnLFxyXG4gICAgICAgICdlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9hdmF0YXIvc3R5bGUvY3NzJyxcclxuICAgICAgICAnZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvZGVzY3JpcHRpb25zL3N0eWxlL2NzcycsXHJcbiAgICAgICAgJ2VsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL2Rlc2NyaXB0aW9ucy1pdGVtL3N0eWxlL2NzcycsXHJcbiAgICAgICAgJ2VsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL2NoZWNrYm94LWdyb3VwL3N0eWxlL2NzcycsXHJcbiAgICAgICAgJ2VsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL3Byb2dyZXNzL3N0eWxlL2NzcycsXHJcbiAgICAgICAgJ2VsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL2ltYWdlLXZpZXdlci9zdHlsZS9jc3MnLFxyXG4gICAgICAgICdlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9lbXB0eS9zdHlsZS9jc3MnLFxyXG4gICAgICAgICdlbGVtZW50LXBsdXMvZXMvY29tcG9uZW50cy9zZWdtZW50ZWQvc3R5bGUvY3NzJyxcclxuICAgICAgICAnZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvY2FsZW5kYXIvc3R5bGUvY3NzJyxcclxuICAgICAgICAnZWxlbWVudC1wbHVzL2VzL2NvbXBvbmVudHMvbWVzc2FnZS9zdHlsZS9jc3MnLFxyXG4gICAgICAgICd4bHN4JyxcclxuICAgICAgICAnZmlsZS1zYXZlcicsXHJcbiAgICAgICAgJ2VsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL3RpbWVsaW5lL3N0eWxlL2NzcycsXHJcbiAgICAgICAgJ2VsZW1lbnQtcGx1cy9lcy9jb21wb25lbnRzL3RpbWVsaW5lLWl0ZW0vc3R5bGUvY3NzJyxcclxuICAgICAgICAndnVlLWltZy1jdXR0ZXInXHJcbiAgICAgIF1cclxuICAgIH0sXHJcbiAgICBjc3M6IHtcclxuICAgICAgcHJlcHJvY2Vzc29yT3B0aW9uczoge1xyXG4gICAgICAgIC8vIHNhc3MgdmFyaWFibGUgYW5kIG1peGluXHJcbiAgICAgICAgc2Nzczoge1xyXG4gICAgICAgICAgYXBpOiAnbW9kZXJuLWNvbXBpbGVyJyxcclxuICAgICAgICAgIGFkZGl0aW9uYWxEYXRhOiBgXHJcbiAgICAgICAgICAgIEB1c2UgXCJAc3R5bGVzL3ZhcmlhYmxlcy5zY3NzXCIgYXMgKjsgQHVzZSBcIkBzdHlsZXMvbWl4aW4uc2Nzc1wiIGFzICo7XHJcbiAgICAgICAgICBgXHJcbiAgICAgICAgfVxyXG4gICAgICB9LFxyXG4gICAgICBwb3N0Y3NzOiB7XHJcbiAgICAgICAgcGx1Z2luczogW1xyXG4gICAgICAgICAge1xyXG4gICAgICAgICAgICBwb3N0Y3NzUGx1Z2luOiAnaW50ZXJuYWw6Y2hhcnNldC1yZW1vdmFsJyxcclxuICAgICAgICAgICAgQXRSdWxlOiB7XHJcbiAgICAgICAgICAgICAgY2hhcnNldDogKGF0UnVsZSkgPT4ge1xyXG4gICAgICAgICAgICAgICAgaWYgKGF0UnVsZS5uYW1lID09PSAnY2hhcnNldCcpIHtcclxuICAgICAgICAgICAgICAgICAgYXRSdWxlLnJlbW92ZSgpXHJcbiAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgICB9XHJcbiAgICAgICAgXVxyXG4gICAgICB9XHJcbiAgICB9XHJcbiAgfSlcclxufVxyXG5cclxuZnVuY3Rpb24gcmVzb2x2ZVBhdGgocGF0aHM6IHN0cmluZykge1xyXG4gIHJldHVybiBwYXRoLnJlc29sdmUoX19kaXJuYW1lLCBwYXRocylcclxufVxyXG4iXSwKICAibWFwcGluZ3MiOiAiO0FBQTJRLFNBQVMsY0FBYyxlQUFlO0FBQ2pULE9BQU8sU0FBUztBQUNoQixPQUFPLFVBQVU7QUFDakIsT0FBTyxxQkFBcUI7QUFDNUIsT0FBTyxnQkFBZ0I7QUFDdkIsT0FBTyxnQkFBZ0I7QUFDdkIsU0FBUywyQkFBMkI7QUFDcEMsU0FBUyxxQkFBcUI7QUFLOUIsT0FBTyxpQkFBaUI7QUFaeEIsSUFBTSxtQ0FBbUM7QUFBNEgsSUFBTSwyQ0FBMkM7QUFjdE4sSUFBTyxzQkFBUSxDQUFDLEVBQUUsS0FBSyxNQUF3QjtBQUM3QyxRQUFNLE9BQU8sUUFB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O0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxNQW9DRCxZQUFZO0FBQUEsSUFDZDtBQUFBO0FBQUEsSUFFQSxjQUFjO0FBQUEsTUFDWixTQUFTO0FBQUEsUUFDUDtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsTUFDRjtBQUFBLElBQ0Y7QUFBQSxJQUNBLEtBQUs7QUFBQSxNQUNILHFCQUFxQjtBQUFBO0FBQUEsUUFFbkIsTUFBTTtBQUFBLFVBQ0osS0FBSztBQUFBLFVBQ0wsZ0JBQWdCO0FBQUE7QUFBQTtBQUFBLFFBR2xCO0FBQUEsTUFDRjtBQUFBLE1BQ0EsU0FBUztBQUFBLFFBQ1AsU0FBUztBQUFBLFVBQ1A7QUFBQSxZQUNFLGVBQWU7QUFBQSxZQUNmLFFBQVE7QUFBQSxjQUNOLFNBQVMsQ0FBQyxXQUFXO0FBQ25CLG9CQUFJLE9BQU8sU0FBUyxXQUFXO0FBQzdCLHlCQUFPLE9BQU87QUFBQSxnQkFDaEI7QUFBQSxjQUNGO0FBQUEsWUFDRjtBQUFBLFVBQ0Y7QUFBQSxRQUNGO0FBQUEsTUFDRjtBQUFBLElBQ0Y7QUFBQSxFQUNGLENBQUM7QUFDSDtBQUVBLFNBQVMsWUFBWSxPQUFlO0FBQ2xDLFNBQU8sS0FBSyxRQUFRLGtDQUFXLEtBQUs7QUFDdEM7IiwKICAibmFtZXMiOiBbInBhdGgiXQp9Cg==
