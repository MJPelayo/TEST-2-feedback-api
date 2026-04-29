// File: script.js
// Purpose: Contains all client-side JavaScript logic for API interaction

// ========== CONFIGURATION ==========
// API endpoint - change if running on different port/host
const API_URL = 'http://localhost:8080';

// ========== DOM ELEMENTS ==========
const feedbackForm = document.getElementById('feedbackForm');
const formMessage = document.getElementById('formMessage');
const feedbackContainer = document.getElementById('feedbackContainer');
const refreshBtn = document.getElementById('refreshBtn');

// ========== API CALL FUNCTIONS ==========

/**
 * Submits feedback to the API
 * @param {Object} formData - The feedback data to submit
 * @returns {Promise} - API response
 */
async function submitFeedbackToAPI(formData) {
    const response = await fetch(`${API_URL}/api/feedback`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData)
    });
    
    const data = await response.json();
    return { response, data };
}

/**
 * Fetches all feedback from the API
 * @returns {Promise} - Array of feedback entries
 */
async function fetchAllFeedback() {
    const response = await fetch(`${API_URL}/api/feedback/all`);
    const data = await response.json();
    return { response, data };
}

// ========== UI HELPER FUNCTIONS ==========

/**
 * Shows a message to the user
 * @param {string} type - 'success' or 'error'
 * @param {string} text - Message text to display
 */
function showMessage(type, text) {
    formMessage.className = `message ${type}`;
    formMessage.textContent = text;
    
    // Auto-hide success messages after 3 seconds
    if (type === 'success') {
        setTimeout(() => {
            formMessage.className = 'message';
        }, 3000);
    }
}

/**
 * Clears the feedback form inputs
 */
function clearForm() {
    feedbackForm.reset();
}

/**
 * Escapes HTML to prevent XSS attacks
 * @param {string} str - String to escape
 * @returns {string} - Escaped string
 */
function escapeHtml(str) {
    if (!str) return '';
    return str
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;');
}

/**
 * Formats a date for display
 * @param {string} dateString - ISO date string
 * @returns {string} - Formatted date
 */
function formatDate(dateString) {
    return new Date(dateString).toLocaleString();
}

/**
 * Renders feedback entries to the container
 * @param {Array} feedbacks - Array of feedback objects
 */
function renderFeedback(feedbacks) {
    if (feedbacks.length === 0) {
        feedbackContainer.innerHTML = '<div class="loading">No feedback submitted yet. Be the first!</div>';
        return;
    }
    
    // Display each feedback entry
    feedbackContainer.innerHTML = feedbacks.map(feedback => `
        <div class="feedback-item">
            <div class="feedback-header">
                <span class="feedback-name">👤 ${escapeHtml(feedback.name)}</span>
                <span class="feedback-email">📧 ${escapeHtml(feedback.email)}</span>
            </div>
            <div class="feedback-subject">📌 ${escapeHtml(feedback.subject)}</div>
            <div class="feedback-message">${escapeHtml(feedback.message)}</div>
            <div class="feedback-date">📅 ${formatDate(feedback.created_at)}</div>
        </div>
    `).join('');
}

/**
 * Shows loading state in feedback container
 */
function showLoading() {
    feedbackContainer.innerHTML = '<div class="loading">Loading feedback...</div>';
}

/**
 * Shows error state in feedback container
 */
function showError(message) {
    feedbackContainer.innerHTML = `<div class="loading">❌ ${message}</div>`;
}

// ========== MAIN APPLICATION LOGIC ==========

/**
 * Loads and displays all feedback
 */
async function loadFeedback() {
    showLoading();
    
    try {
        const { response, data } = await fetchAllFeedback();
        
        if (response.ok) {
            renderFeedback(data);
        } else {
            showError('Failed to load feedback. Server error.');
        }
    } catch (error) {
        console.error('Error loading feedback:', error);
        showError('Failed to connect to server. Make sure the API is running on ' + API_URL);
    }
}

/**
 * Handles form submission
 * @param {Event} event - Form submit event
 */
async function handleFormSubmit(event) {
    event.preventDefault();
    
    // Get form data
    const formData = {
        name: document.getElementById('name').value,
        email: document.getElementById('email').value,
        subject: document.getElementById('subject').value,
        message: document.getElementById('message').value
    };
    
    try {
        const { response, data } = await submitFeedbackToAPI(formData);
        
        if (response.ok) {
            // Success
            showMessage('success', '✅ ' + data.message);
            clearForm();
            // Refresh the feedback list
            await loadFeedback();
        } else {
            // Error from server
            showMessage('error', '❌ Error: ' + data.error);
        }
    } catch (error) {
        console.error('Error submitting feedback:', error);
        showMessage('error', '❌ Could not connect to server. Make sure the API is running on ' + API_URL);
    }
}

// ========== EVENT LISTENERS ==========
feedbackForm.addEventListener('submit', handleFormSubmit);
refreshBtn.addEventListener('click', loadFeedback);

// ========== INITIALIZATION ==========
// Load feedback when page loads
loadFeedback();