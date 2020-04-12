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
        submit: 'Submit',
        show_more: 'Load more comments',
        show_more_depth: 'Load more replies',
        format_date: (date) => `${date.getFullYear()}.${date.getMonth()}.${date.getDate()} ${date.getHours()}:${date.getMinutes()}`
    };

    // Optional options
    const batchSize = options.batchSize || 10;
    const maxDepth = options.maxDepth || 2;

    // Persistent values
    window._gomment = {
        comments: [],
        lastOffset: 0
    };

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
    const commentsElement = insertElement('div', 'gomment-comments', container);

    function queryComments(thread, offset, max, depth) {
        const apiURL = api.endsWith('/') ? api : `/${api}`;
        return window.fetch(`${apiURL}comments?thread=${thread}&offset=${offset}&max=${max}&depth=${depth}`);
    }

    // Render single comment and children recursively
    function renderComment(parent, comment, treeIndex, depth) {
        const commentElement = insertElement('div', 'gomment-comment', parent);
        insertElement('div', 'gomment-comment-author', commentElement, { innerHTML: comment.comment.author });
        insertElement('div', 'gomment-comment-date', commentElement, { innerHTML: i18n.format_date(new Date(comment.comment.created_at)) });
        insertElement('div', 'gomment-comment-text', commentElement, { innerHTML: comment.comment.text });
        const childrenElement = insertElement('div', 'gomment-comment-children', commentElement);

        if (comment.children) {
            for (const childComment of comment.children) {
                renderComment(childrenElement, childComment, treeIndex, depth + 1);
            }
        } else if (comment.has_children) {
            // No children but hasChildren -> Load more button
            insertElement('button', 'gomment-show-more-depth-button', childrenElement, {
                innerHTML: i18n.show_more_depth,
                onclick: () => loadMoreDepth(treeIndex, depth)
            });
        }
    }

    // Render all (top level) comments with children
    function renderComments(comments, total) {
        // Remove all comments from element
        commentsElement.innerHTML = '';
        // Render comments
        comments.forEach((comment, index) => renderComment(commentsElement, comment, index, 0));
        // Show more button
        if (comments.length < total) {
            const showMoreContainer = insertElement('div', 'gomment-show-more-container', commentsElement);
            insertElement('button', 'gomment-show-more-button', showMoreContainer, {
                innerHTML: i18n.show_more,
                onclick: () => loadNextBatch()
            });
        }
    }

    // Load and render comments
    function loadComments(offset, max, depth) {
        const handler = (data) => {
            data.comments.forEach((item, index) => window._gomment.comments[index + offset] = item);
            window._gomment.lastOffset = offset;
            renderComments(window._gomment.comments, data.total);
        };

        queryComments(thread, offset, max, depth)
            .then(data => data.json())
            .then(handler);

    }

    // Load intial comments
    function loadCommentsInitial() {
        loadComments(0, batchSize, maxDepth);
    }

    // Load 'more' comments
    function loadNextBatch() {
        loadComments(window._gomment.lastOffset + batchSize, batchSize, maxDepth);
    }

    // Load more replies
    function loadMoreDepth(index, depth) {
        loadComments(index, 1, depth + batchSize);
    }

    loadCommentsInitial();
}
