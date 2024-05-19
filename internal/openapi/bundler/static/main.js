/**
 * Configuration object for a URL in SwaggerUIConfig.
 * @typedef {Object} UrlConfig
 * @property {string} url - The URL pointing to the Swagger/OpenAPI definition.
 * @property {string} name - The name associated with the URL.
 */

/**
 *
 * Represents the configuration options for SwaggerUI.
 * @typedef {Object} SwaggerUIConfig
 * @property {Object} [plugins] - Optional plugins for SwaggerUI.
 * @property {string} [dom_id] - ID of the DOM element to render SwaggerUI.
 * @property {Object} [presets] - Presets configuration for SwaggerUI.
 * @property {Object} [layout] - Layout configuration for SwaggerUI.
 * @property {boolean} [deepLinking] - Enable deep linking for SwaggerUI.
 * @property {Object} [requestInterceptor] - Interceptor for outgoing requests in SwaggerUI.
 * @property {Object} [responseInterceptor] - Interceptor for incoming responses in SwaggerUI.
 * @property {string} [oauth2RedirectUrl] - OAuth2 redirect URL for SwaggerUI.
 * @property {string} [url] - URL to Swagger API definition.
 * @property {UrlConfig[]} [urls] - URLs configuration for SwaggerUI, represented as an array of strings.
 * @property {string} [validatorUrl] - Validator URL for SwaggerUI.
 * @property {Object} [spec] - Swagger API specification for SwaggerUI.
 * @property {Object} [specActions] - Specification actions for SwaggerUI.
 * @property {Object} [filter] - Filter configuration for SwaggerUI.
 * @property {Object} [docExpansion] - Document expansion configuration for SwaggerUI.
 * @property {Object} [operationsSorter] - Operations sorter configuration for SwaggerUI.
 * @property {boolean} [showMutatedRequest] - Show mutated request configuration for SwaggerUI.
 * @property {string[]} [supportedSubmitMethods] - Supported submit methods for SwaggerUI, represented as an array of strings.
 * @property {boolean} [modelPropertyMacro] - Model property macro configuration for SwaggerUI.
 * @property {boolean} [modelPropertyMacroEnabled] - Enable model property macro for SwaggerUI.
 * @property {boolean} [showRequestHeaders] - Show request headers configuration for SwaggerUI.
 * @property {string} [defaultModelRendering] - Default model rendering configuration for SwaggerUI.
 * @property {boolean} [displayOperationId] - Display operation ID configuration for SwaggerUI.
 * @property {Object} [pluginsConfig] - Plugins configuration for SwaggerUI.
 * @property {Object} [layoutConfig] - Layout configuration for SwaggerUI.
 * @property {Object} [filterConfig] - Filter configuration for SwaggerUI.
 * @property {boolean} [deepLinkingConfig] - Deep linking configuration for SwaggerUI.
 * @property {string} [syntaxHighlight] - Syntax highlight configuration for SwaggerUI.
 * @property {string} [oauth2RedirectUrl] - OAuth2 redirect URL configuration for SwaggerUI.
 * @property {string} [logoUrl] - Logo URL configuration for SwaggerUI.
 */

window.onload = function () {
    /** @type {UrlConfig[]} */
    var urls = JSON.parse(document.getElementById("swaggerUrls").textContent);
    
    /** @type {SwaggerUIConfig} */
    var swaggerConfig = {
      dom_id: "#swagger-ui",
      deepLinking: true,
      presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
      plugins: [SwaggerUIBundle.plugins.DownloadUrl],
      layout: "StandaloneLayout",
      urls: urls,
    };
  window.ui = SwaggerUIBundle(swaggerConfig);
};
