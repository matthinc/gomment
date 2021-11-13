/**
 * Content of a single comment.
 * @typedef {Object} CommentModel
 * @property {number} comment_id - Unique identifier for the comment.
 * @property {string} author - The author of the comment.
 * @property {number} created_at - Creation date of the comment.
 * @property {number} touched_at - Touch date of the comment.
 * @property {string} text - The content of the comment.
 * @property {number} num_children - The total number of availabe children, independent of the shown children.
 */

/**
 * A comment response from the backend.
 * @typedef {Object} CommentTreeNode
 * @property {CommentModel} comment - The actual comment data (content).
 * @property {Array<CommentTreeNode>} children - Children comments of this comment.
 * @property {CommentDom | null} dom - The DOM node the comment was rendered into, null otherwise.
 */

/**
 * A response returned by the backend for a comment query.
 * @typedef {Object} CommentQueryResponse
 * @property {Array<CommentTreeNode>} comments - The queried comments.
 * @property {number} num_total - The total amount of comments available in the thread.
 * @property {number} num_root - The amount of root comments.
 * @property {number} num_root_payload - The amount of comments returned by the query.
 * @property {number} thread_id - The identifier of this thread.
 */

/**
 * An element representing a DOM node for a comment, containing references to the most important elements.
 * @typedef {Object} CommentDom
 * @property {HTMLElement} elRoot - .
 * @property {HTMLElement | null} elAuthor - .
 * @property {HTMLElement | null} elDate - .
 * @property {HTMLElement | null} elText - .
 * @property {HTMLElement | null} elReply - .
 * @property {HTMLElement | null} elMoreSiblings - .
 * @property {HTMLElement | null} elMoreChildren - .
 * @property {HTMLElement} elChildren - .
 */

/**
 * A collection of all relevent DOM elements in the comment creation dialog.
 * @typedef {Object} InputSectionDOM
 * @property {HTMLElement} elRoot - .
 * @property {HTMLInputElement} elMail - .
 * @property {HTMLInputElement} elName - .
 * @property {HTMLTextAreaElement} elContent - .
 * @property {HTMLSpanElement} elError - .
 * @property {HTMLButtonElement} elSubmit - .
 */

/**
 * Creates a new element and inserts it as the child of parent.
 * @param {string} type - Type of the new DOM element.
 * @param {string} className - Name of the class to assign to the newly created element.
 * @param {HTMLElement | null} parent - The new parent of the DOM element, null if it should not be attached to a parent.
 * @param {*} attributes - Attributes to assign to the new DOM element.
 * @returns {HTMLElement} The newly created child.
 */
function insertElement(type, className, parent, attributes = {}) {
  if (parent === undefined) {
    throw new Error('No valid parent element was provided.');
  }

  const elem = document.createElement(type);
  elem.className = className;
  // Misc. attributes
  for (const key in attributes) {
    // @ts-ignore
    elem[key] = attributes[key];
  }
  if(parent) {
    parent.appendChild(elem);
  }
  return elem;
}

/**
* Transforms a date to a string that will be shown next to the username
* @param {Date} date - the date
* @returns {string}
*/
function defaultDateTransformer(date) {
  return date.toLocaleString();
}

export class Gomment {

  /**
   * @constructor
   * @param {?*} [options]
   */
  constructor(options) {
    if (!options || typeof options !== 'object') {
      throw new Error('mandatory parameter `options` was not provided');
    }

    // Required options
    /** @type {string} */
    this.apiURL = options.apiURL;
    // append a missing trailing backslash
    this.apiURL = this.apiURL.endsWith('/') ? this.apiURL : `${this.apiURL}/`;

    /** @type {string} */
    if(typeof options.thread === 'string') {
      this.threadPath = options.threadPath;
    } else {
      this.threadPath = this.getThreadPathFromLocation();
    }

    // Optional options
    /** @type {number} */
    this.batchSize = options.batchSize || 10;
    /** @type {number} */
    this.maxDepth = options.maxDepth || 6;

    // i18n
    this.i18n = options.i18n || {
      title: 'Comments',
      empty: 'No comments',
      placeholder_name: 'Name',
      placeholder_email: 'E-Mail',
      placeholder_text: 'Your comment',
      submit: 'Submit',
      submit_reply: 'Reply',
      show_more: 'Load more comments',
      show_more_depth: 'Load more replies',
      alert_missing_information: 'Please fill out all required fields!',
      reply: 'Reply',
      new_comment: 'Write comment',
      /** @type {(date: Date) => string} */
      format_date: defaultDateTransformer
    };

    // comment input section
    /** @type {HTMLButtonElement | null} */
    this.newCommentButton = null;

    /** @type {InputSectionDOM | null} */
    this.inputSectionDom = null;

    /** @type {number} */
    this.replyRecipient = 0;

    /** @type {HTMLElement | null} */
    this.replyIndicator = null;

    // stateful information
    /** @type {number | null} */
    this.threadId = null;

    /** @type {HTMLElement | null} */
    this.newButton = null;

    /** @type {CommentTreeNode} */
    this.rootNode = {
      comment: {
        num_children: 0,
        comment_id: 0,
        created_at: 0,
        touched_at: 0,
        author: '',
        text: '',
      },
      children: [],
      dom: null,
    };

    /** @type {number} */
    this.numTotal = 0;
  }

  /**
   * Set various thread attributes and update the DOM accordingly.
   * @param {CommentQueryResponse} queryResponse - The response for querying comments initially.
   */
  setThreadMetadata(queryResponse) {
    this.numTotal = queryResponse.num_total;
    this.rootNode.comment.num_children = queryResponse.num_root;
    this.threadId = queryResponse.thread_id;
    console.warn("TODO: update DOM in setTotalComments");
  }

  /**
   * Query for comments with the specified parameters and render them.
   * @param {number} max - Maximum amount of comments to query.
   * @param {number} depth - Maximum level of depth to query.
   * @returns {void}
   */
  loadComments(max, depth) {
    window.fetch(`${this.apiURL}comments/nbf?threadPath=${encodeURIComponent(this.threadPath)}&max=${max}&depth=${depth}`)
      .then(rawData => rawData.json())
      .then(/** @type {function(CommentQueryResponse): void} */ jsonData => {
        this.setThreadMetadata(jsonData);
        this.rootNode.children = jsonData.comments;

        this.renderComment(this.rootNode);
      });
  }

  /**
   * Load comments with the initial default parameters.
   * @returns {void}
   */
  loadCommentsInitial() {
    this.loadComments(this.batchSize, this.maxDepth);
  }

  /**
   * Load 'more' sibling comments.
   * @param {CommentTreeNode} parent - The parent for which to load more comments.
   * @returns {void}
   */
  loadMoreSiblings(parent) {
    const childComments = parent.children;
    // order the id's ascending - as required by the API
    const excludeIds = childComments.map(c => c.comment.comment_id).sort((a, b) => a - b).join(',');

    // only load comments older than the newest comment. if no comment
    // is present take an arbitrary high number
    const newestCreatedAt = childComments.reduce((previous, current) => {
      return Math.max(previous, current.comment.touched_at);
    }, 0) || 0x7FFFFFFFFFFFF;

    window
      .fetch(`${this.apiURL}morecomments/nbf?threadId=${this.threadId}&parentId=${parent.comment.comment_id}&newestCreatedAt=${newestCreatedAt}&limit=${this.batchSize}&excludeIds=${excludeIds}`)
      .then(rawData => rawData.json())
      .then(/** @type {function(Array<CommentModel>): void} */ comments => {
        comments.forEach(c => {
          /** @type {CommentTreeNode} */
          const treeNode = {
            comment: c,
            children: [],
            dom: null,
          };
          parent.children.push(treeNode);
        });
        this.renderComment(parent);
      });
  }

  /**
   * Create a DOM node for displaying a comment and return references to specific elements.
   * @param {boolean} isRootComment - indicates whether content-specific DOM element shall be omited.
   * @returns {CommentDom}
   */
  createEmptyCommentDom(isRootComment) {
    const elRoot = insertElement('div', isRootComment ? 'gomment-comments' : 'gomment-comment', null);

    let elAuthor = null;
    let elDate = null
    let elText = null;
    let elReply = null;
    if (!isRootComment) {
      elAuthor = insertElement('div', 'gomment-comment-author', elRoot);
      elDate = insertElement('div', 'gomment-comment-date', elRoot);
      elText = insertElement('div', 'gomment-comment-text', elRoot);
      elReply = insertElement('a', 'gomment-comment-reply', elRoot, { innerHTML: this.i18n.reply });
    }
    const elChildren = insertElement('div', 'gomment-comment-children', elRoot);

    return {
      elRoot,
      elAuthor,
      elDate,
      elText,
      elReply,
      elMoreSiblings: null,
      elMoreChildren: null,
      elChildren,
    };
  }

  /**
   * Render single comment and children recursively.
   * @param {CommentTreeNode} parentNode - The new comment to render.
   * @returns {void}
   */
  renderComment(parentNode) {
    const parentDom = parentNode.dom;
    if (!parentDom) {
      throw new Error('failed precondition: parent DOM must be created before rendering child comment');
    }

    parentNode.children.forEach(childNode => {
      // skip if the child was already rendered
      if(childNode.dom) {
        return;
      }

      const dom = this.createEmptyCommentDom(false);
      if (!dom.elAuthor || !dom.elDate || !dom.elText || !dom.elReply) {
        throw new Error('failed to create comment DOM element');
      }

      dom.elAuthor.innerHTML = childNode.comment.author;
      dom.elDate.innerHTML = this.i18n.format_date(new Date(childNode.comment.created_at * 1000));
      dom.elText.innerHTML = childNode.comment.text;

      dom.elReply.onclick = () => {
        this.moveInputSection(dom.elChildren, childNode.comment.comment_id);
      };

      // attach child to the parent dom
      childNode.dom = dom;
      parentDom.elChildren.appendChild(dom.elRoot);

      // recurse over child
      this.renderComment(childNode);
    });

    // "show more" button
    if (parentNode.children.length > 0) {
      if (!parentDom.elMoreSiblings && parentNode.comment.num_children > parentNode.children.length) {
        parentDom.elMoreSiblings = insertElement('div', 'gomment-show-more-container', parentDom.elRoot);
        insertElement('button', 'gomment-show-more-button', parentDom.elMoreSiblings, {
          innerHTML: this.i18n.show_more,
          onclick: () => this.loadMoreSiblings(parentNode),
        });
      } else if(parentDom.elMoreSiblings) {
        const el = parentDom.elMoreSiblings;
        el.parentElement && el.parentElement.removeChild(el);
        parentDom.elMoreSiblings = null;
      }
    }

    // "load children" button
    if (!parentDom.elMoreChildren && parentNode.children.length === 0 && parentNode.comment.num_children > 0) {
      // No children but hasChildren -> Load more button
      parentDom.elMoreChildren = insertElement('button', 'gomment-show-more-depth-button', parentDom.elChildren, {
        innerHTML: this.i18n.show_more_depth,
        onclick: () => this.loadMoreSiblings(parentNode)
      });
    } else if (parentDom.elMoreChildren) {
      const el = parentDom.elMoreChildren
      el.parentElement && el.parentElement.removeChild(el);
      parentDom.elMoreChildren = null;
    }
  }

  /**
   * Publish new comment to this thread
   * @param {number} parent - parent comment (0 = top-level)
   * @param {string} name
   * @param {string} email
   * @param {string} content - text content
   * @returns {Promise<void>} true if data was valid
   */
  publishComment(parent, name, email, content) {
    return new Promise((resolve, reject) => {
      if (!name || !email || !content) {
        return reject(this.i18n.alert_missing_information);
      }

      const data = {
        author: name,
        email,
        text: content,
        parent_id: parent,
        thread_path: this.threadPath,
      };

      // Post comment
      window.fetch(`${this.apiURL}comment`, {
        method: 'POST',
        body: JSON.stringify(data)
      })
        .then(response => response.json())
        .then(data => {
          this.newComment = data.id;
          return resolve();
        })
        .catch(err => {
          return reject(err);
        });
    });
  }

  /**
   * Move comment input to new parent and change recipient
   * @param {HTMLElement} newParent - new parent
   * @param {number} recipient - new recipient (0 for top-level comments)
   */
  moveInputSection(newParent, recipient) {
    if(this.inputSectionDom === null || this.newCommentButton === null) {
      throw new Error('precondition failed: inputSectionDom and newCommentButton must be created before calling moveInputSection');
    }

    if (recipient < 0) {
      throw new Error('reply recipient must be a valid comment_id');
    }

    // Hide new comment button
    this.newCommentButton.hidden = recipient === 0;

    // Set recipient
    this.replyRecipient = recipient;

    // Move
    if (newParent.childNodes.length === 0) {
      newParent.appendChild(this.inputSectionDom.elRoot);
    } else {
      newParent.insertBefore(this.inputSectionDom.elRoot, newParent.childNodes[0]);
    }
  }

  /**
   * Hide the comment input section, but make the root level comment button visible.
   * @returns {void}
   */
  hideInputSection() {
    if(this.inputSectionDom === null || this.newCommentButton === null) {
      throw new Error('precondition failed: inputSectionDom and newCommentButton must be created before calling moveInputSection');
    }

    // show the 'new comment' button
    this.newCommentButton.hidden = false;

    /** @type {HTMLElement | null} */
    const currentParent = this.inputSectionDom.elRoot.parentElement;
    currentParent && currentParent.removeChild(this.inputSectionDom.elRoot);
  }

  /**
   * Create the comment input section DOM elements.
   * @returns {InputSectionDOM}
   */
  createInputSection() {
    const elRoot = insertElement('div', 'gomment-input-section', null);
    const elMail = /** @type {HTMLInputElement} */ (insertElement('input', 'gomment-email', elRoot, { placeholder: this.i18n.placeholder_email }));
    const elName = /** @type {HTMLInputElement} */ (insertElement('input', 'gomment-display-name', elRoot, { placeholder: this.i18n.placeholder_name }));
    const elContent = /** @type {HTMLTextAreaElement} */ (insertElement('textarea', 'gomment-text-input', elRoot, { placeholder: this.i18n.placeholder_text }));
    const elError = /** @type {HTMLInputElement} */ (insertElement('span', 'gomment-new-error', elRoot));
    const elSubmit = /** @type {HTMLButtonElement} */ (insertElement('button', 'gomment-submit-button', elRoot, { innerHTML: this.i18n.submit}));

    elSubmit.addEventListener('click', this.onSendComment.bind(this));

    return {
      elRoot,
      elMail,
      elName,
      elContent,
      elError,
      elSubmit,
    };
  }

  onSendComment() {
    if(this.inputSectionDom === null) {
      throw new Error('precondition failed: onSendComment requires inputSectionDom to be non-null');
    }

    /** @type {InputSectionDOM} */
    const d = this.inputSectionDom;

    // disable all input elements to give visual indication
    const elements = [
      d.elMail,
      d.elName,
      d.elContent,
      d.elSubmit,
    ];
    elements.forEach(e => e.disabled = true);

    this.publishComment(
      this.replyRecipient,
      d.elName.value,
      d.elMail.value,
      d.elContent.value,
    )
      .then(() => {
        // Clear inputs
        d.elName.value = '';
        d.elMail.value = '';
        d.elContent.value = '';
        d.elError.innerText = '';

        this.hideInputSection();

        // enable all elements in the end
        elements.forEach(e => e.disabled = false);
      })
      .catch(err => {
        d.elError.innerText = err.toString();

        // enable all elements in the end
        elements.forEach(e => e.disabled = false);
      });
  }

  /**
   * Injects the gomment instance into a parent element.
   * @param {string | HTMLElement} element - Parent element.
   */
  injectInto(element) {
    /** @type {HTMLElement} */
    let container;
    if (typeof element === 'string') {
      const el = document.querySelector(element);
      if (!el || !(el instanceof HTMLElement)) {
        throw new Error('HTML element with the specifier "${element}" was not found.');
      }
      container = insertElement('div', 'gomment', el);
    } else if (element instanceof HTMLElement) {
      container = element;
    } else {
      throw new Error('Parent element needs to be provided as a query selector or `HTMLElement`');
    }

    // Create container element
    insertElement('div', 'gomment-title', container, { innerHTML: this.i18n.title });

    // create container at the top of the comments for the input section
    const topInputSectionContainer = insertElement('div', 'gomment-new-comment-input-container', container);

    // create button for moving the comment section to the top level
    const newButton =  /** @type {HTMLButtonElement} */ (insertElement('button', 'gomment-new-comment', topInputSectionContainer, { innerHTML: this.i18n.new_comment}));
    newButton.addEventListener('click', e => {
      this.moveInputSection(topInputSectionContainer, 0);
    })
    this.newCommentButton = newButton

    this.inputSectionDom = this.createInputSection();
    topInputSectionContainer.appendChild(this.inputSectionDom.elRoot);

    // Initial input field position
    this.moveInputSection(topInputSectionContainer, 0);

    // Comments section
    this.rootNode.dom = this.createEmptyCommentDom(true);
    container.appendChild(this.rootNode.dom.elRoot);

    // Load and render comments
    this.loadCommentsInitial();
  }

  /**
   * Calculate a cleaned-up thread path based on the current browser
   * location by trimming all trailing slashes.
   * @returns {string} - The cleaned up thread path
   */
  getThreadPathFromLocation() {
    let pathname = window.location.pathname;

    // remove trailing slash if it exists
    if (pathname[pathname.length - 1] === '/') {
      pathname = pathname.substring(0, pathname.length - 1);
    }

    // if the path is empty, use a slash as path
    if (pathname.length === 0) {
      pathname = '/';
    }

    return pathname;
  }
}
