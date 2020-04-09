// Define `gomment` function
window.gomment = function (options) {
    // Required options
    const { element, api, thread } = options;

    // i18n
    const i18n = options.i18n || {
        title: 'Comments',
        empty: 'No comments',
        input_title: 'Leave a comment:',
        placeholder_name: 'Name',
        placeholder_email: 'E-Mail',
        placeholder_text: 'Your comment',
        submit: 'Submit'
    };

    // Optional options
    const batchSize = options.batchSize || 10;
    const maxDepth = options.maxDepth || 2;

    function insertElement(type, className, parent, attributes = {}) {
        const elem = document.createElement(type);
        elem.className = className;
        // Misc. attributes
        for (const key in attributes) {
            elem[key] = attributes[key];
        }
        parent.appendChild(elem);
        return elem;
    }

    // Create container element
    const container = insertElement('div', 'gomment', document.getElementById(element));
    insertElement('div', 'gomment-title', container, { innerHTML: i18n.title });

    // Comment input section
    const inputSection = insertElement('div', 'gomment-input-section', container);
    insertElement('div', 'gomment-input-title', inputSection, { innerHTML: i18n.input_title });
    insertElement('input', 'gomment-email', inputSection, { placeholder: i18n.placeholder_email });
    insertElement('input', 'gomment-display-name', inputSection, { placeholder: i18n.placeholder_name });
    insertElement('textarea', 'gomment-text-input', inputSection, { placeholder: i18n.placeholder_text });
    insertElement('button', 'gomment-submit-button', inputSection, { innerHTML: i18n.submit });

    // Comments section
    const comments = insertElement('div', 'gomment-comments', container);

    function queryComments(thread, offset, max, depth) {
        const apiURL = api.endsWith('/') ? api : `/${api}`;
        return window.fetch(`${apiURL}comments?thread=${thread}&offset=${offset}&max=${max}&depth=${depth}`);
    }

    // Todo: Render comments
    queryComments(thread, 0, batchSize, maxDepth)
        .then(data => data.json())
        .then(data => comments.innerHTML = JSON.stringify(data));
}
