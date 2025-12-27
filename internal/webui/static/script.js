const API_BASE = '/api';

// I18n
const translations = {
    en: {
        login: "Login",
        setup: "Setup Admin Account",
        username: "Username",
        password: "Password",
        newPassword: "New Password",
        newPassword: "New Password",
        servers: "Emby Media Servers",
        addServer: "Add Emby Media Server",
        editServer: "Edit Emby Media Server",
        logs: "Logs",
        users: "User Management",
        configFile: "Config File",
        episodesUnplayPrior: "Episodes unplayed first",
        episodesUnplayPriorDesc: "Reorder episodes to show unplayed ones first",
        resortRandomItems: "Resort random items",
        resortRandomItemsDesc: "Extra randomization to improve list randomness",
        proxyErrorStrategy: "Proxy error strategy",
        proxyErrorStrategyDesc: "origin: back to source; reject: deny",
        imagesQuality: "Images quality",
        imagesQualityDesc: "Range 1–100; recommend 70–90",
        downloadStrategy: "Download strategy",
        downloadStrategyDesc: "403: disable; origin: proxy; direct: redirect",
        cacheEnable: "Enable cache",
        cacheExpired: "Cache expiration",
        cacheWhitelist: "Cache whitelist (regex per line)",
        cacheWhitelistDesc: "Only matched routes are cached; empty uses built-ins",
        vpEnable: "Video preview enable",
        vpContainers: "Video preview containers",
        vpIgnore: "Video preview ignore templates",
        pathEmby2Openlist: "Path emby2openlist",
        logDisableColor: "Disable colored logs",
        strmPathMap: "STRM path-map",
        strmPathMapDesc: "Each line: from => to",
        ltgEnable: "Local tree gen enable",
        ltgFfmpeg: "FFmpeg enable",
        ltgVirtual: "Virtual containers",
        ltgStrm: "STRM containers",
        ltgMusic: "Music containers",
        ltgAutoRemove: "Auto remove max count",
        ltgRefresh: "Refresh interval (minutes)",
        ltgScanPrefixes: "Scan prefixes",
        ltgIgnoreContainers: "Ignore containers",
        ltgThreads: "Threads",
        sslEnable: "Enable HTTPS",
        sslSingle: "Single port",
        sslKey: "SSL Key",
        sslCrt: "SSL Cert",
        logout: "Logout",
        save: "Save",
        cancel: "Cancel",
        clear: "Clear",
        deleteConfirm: "Are you sure you want to delete this server?",
        port: "Go Emby Port",
        mountPath: "Mount Path",
        mountPathDesc: "Multiple paths supported, separate by , or ;",
        localMediaRoot: "Local Media Root",
        localMediaRootDesc: "Files starting with this path will bypass proxy; Multiple paths supported, separate by , or ;",
        dlCacheValue: "Direct Link Cache",
        dlCacheUnit: "Unit",
        serverName: "Server Name",
        internalRedirect: "Enable Internal Redirect",
        disableProxy: "Disable Proxy (Config Only)",
        updatePassword: "Update Password",
        refresh: "Refresh",
        playbackOnly: "Playback Logs",
        errorLogs: "Error Logs",
        allLogs: "All Logs",
        allServers: "all",
        autoFast: "Refresh Interval: 1s",
        autoNormal: "Refresh Interval: 3s",
        time: "Time",
        level: "Level",
        message: "Message",
        networkError: "Network error",
        dbInit: "Initializing...",
        success: "Success",
        setupComplete: "Setup complete. Please login.",
        appTitle: "Go Emby",
        projectAddress: "Project Address (GitHub)"
    },
    zh: {
        login: "登录",
        setup: "设置管理员账户",
        username: "用户名",
        password: "密码",
        newPassword: "新密码",
        newPassword: "新密码",
        servers: "Emby媒体服务器",
        addServer: "添加Emby媒体服务器",
        editServer: "编辑Emby媒体服务器",
        logs: "系统日志",
        users: "用户管理",
        configFile: "配置文件",
        episodesUnplayPrior: "剧集未播优先",
        episodesUnplayPriorDesc: "将未播放的剧集排在前面",
        resortRandomItems: "随机列表二次重排",
        resortRandomItemsDesc: "进一步打乱随机列表提高随机性",
        proxyErrorStrategy: "代理异常策略",
        proxyErrorStrategyDesc: "origin: 回源处理；reject: 拒绝请求",
        imagesQuality: "图片质量",
        imagesQualityDesc: "范围 1–100；建议 70–90",
        downloadStrategy: "下载策略",
        downloadStrategyDesc: "403: 禁用；origin: 代理；direct: 重定向直链",
        cacheEnable: "启用缓存",
        cacheExpired: "缓存过期时间",
        cacheWhitelist: "缓存白名单（每行一个正则）",
        cacheWhitelistDesc: "仅匹配的接口参与缓存；留空使用默认",
        vpEnable: "开启转码资源获取",
        vpContainers: "转码资源容器列表",
        vpIgnore: "忽略转码清晰度",
        pathEmby2Openlist: "挂载路径映射",
        logDisableColor: "禁用彩色日志",
        strmPathMap: "STRM 路径映射",
        strmPathMapDesc: "每行一个映射：from => to",
        ltgEnable: "开启本地目录树生成",
        ltgFfmpeg: "开启 FFmpeg 辅助",
        ltgVirtual: "虚拟容器",
        ltgStrm: "STRM 容器",
        ltgMusic: "音乐容器",
        ltgAutoRemove: "自动删除最大数量",
        ltgRefresh: "刷新间隔（分钟）",
        ltgScanPrefixes: "扫描前缀",
        ltgIgnoreContainers: "忽略容器",
        ltgThreads: "线程数",
        sslEnable: "启用 HTTPS",
        sslSingle: "单一端口",
        sslKey: "私钥文件",
        sslCrt: "证书文件",
        logout: "退出登录",
        save: "保存",
        cancel: "取消",
        clear: "清空",
        deleteConfirm: "确定要删除此服务器吗？",
        port: "Go Emby 端口",
        mountPath: "挂载路径",
        mountPathDesc: "支持多个路径，使用逗号或分号分隔",
        localMediaRoot: "本地媒体根路径",
        localMediaRootDesc: "以此路径开头的文件将绕过代理直连播放；支持多个路径，使用逗号或分号分隔",
        dlCacheValue: "直链缓存时间",
        dlCacheUnit: "单位",
        serverName: "服务器名称",
        internalRedirect: "开启内部重定向",
        disableProxy: "禁用代理 (仅保留配置)",
        updatePassword: "修改密码",
        refresh: "刷新",
        playbackOnly: "播放日志",
        errorLogs: "错误日志",
        allLogs: "全部日志",
        allServers: "all",
        autoFast: "刷新时间：1秒",
        autoNormal: "刷新时间：3秒",
        time: "时间",
        level: "级别",
        message: "消息",
        networkError: "网络错误",
        dbInit: "初始化中...",
        success: "成功",
        setupComplete: "设置完成，请登录。",
        appTitle: "Go Emby",
        projectAddress: "项目地址 (GitHub)"
    }
};

const storedLang = localStorage.getItem('lang');
let lang = storedLang || ((navigator.language || navigator.userLanguage).startsWith('zh') ? 'zh' : 'en');
const t = (key) => translations[lang][key] || key;

window.toggleLang = () => {
    lang = lang === 'zh' ? 'en' : 'zh';
    localStorage.setItem('lang', lang);
    applyTranslations();
    checkInit();
};

function updateLangButtons() {
    const text = lang === 'zh' ? 'English' : '简体中文';
    const authBtn = document.getElementById('lang-btn-auth');
    const sidebarBtn = document.getElementById('lang-btn-sidebar');
    if (authBtn) authBtn.textContent = text;
    if (sidebarBtn) sidebarBtn.textContent = text;
}

// State
let authToken = localStorage.getItem('token');
let servers = [];

// DOM Elements
const authScreen = document.getElementById('auth-screen');
const dashboard = document.getElementById('dashboard');
const authForm = document.getElementById('auth-form');
const authTitle = document.getElementById('auth-title');
const authSubtitle = document.getElementById('auth-subtitle');
const authError = document.getElementById('auth-error');
const logoutBtn = document.getElementById('logout-btn');
const menuItems = document.querySelectorAll('.menu li');
const pages = document.querySelectorAll('.page');

// Init
async function init() {
    applyTranslations();
    checkInit();
    if (authToken) {
        showDashboard();
    } else {
        showAuth();
    }
}

function applyTranslations() {
    document.querySelectorAll('[data-t]').forEach(el => {
        const key = el.dataset.t;
        if (el.tagName === 'INPUT' && el.placeholder) {
            el.placeholder = t(key);
        } else {
            el.textContent = t(key);
        }
    });
    updateLangButtons();
}

async function checkInit() {
    try {
        const res = await fetch(`${API_BASE}/check-init`);
        const data = await res.json();

        // Reset title to appTitle first (it's handled by applyTranslations usually, but good to be safe)
        // But since we use data-t="appTitle", applyTranslations handles it.

        if (authSubtitle) {
            if (!data.initialized) {
                authSubtitle.textContent = t('setup');
                authForm.dataset.mode = "setup";
            } else {
                // For login mode, we don't need a subtitle "Login", just empty is fine or specific text
                authSubtitle.textContent = "";
                authForm.dataset.mode = "login";
            }
        } else {
            // Fallback: if subtitle is missing, append status to title or handle gracefully
            // But we should rely on subtitle existing now
            const titleEl = document.getElementById('auth-title');
            if (!data.initialized) {
                titleEl.textContent = t('setup'); // Only overwrite if setup needed and no subtitle
                authForm.dataset.mode = "setup";
            } else {
                // Do NOT overwrite title for login if we want it to stay "Go Emby"
                // titleEl.textContent = t('login'); 
                authForm.dataset.mode = "login";
            }
        }
    } catch (e) {
        console.error("Init check failed", e);
    }
}

// Auth
authForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    const username = document.getElementById('auth-username').value;
    const password = document.getElementById('auth-password').value;
    const mode = authForm.dataset.mode || 'login';

    try {
        const res = await fetch(`${API_BASE}/${mode}`, {
            method: 'POST',
            body: JSON.stringify({ Username: username, Password: password }),
            headers: { 'Content-Type': 'application/json' }
        });

        if (res.ok) {
            if (mode === 'setup') {
                alert(t('setupComplete'));
                checkInit();
            } else {
                const data = await res.json();
                authToken = data.token;
                localStorage.setItem('token', authToken);
                showDashboard();
            }
        } else {
            const err = await res.json();
            authError.textContent = err.error || 'Authentication failed';
        }
    } catch (e) {
        authError.textContent = t('networkError');
    }
});

logoutBtn.addEventListener('click', () => {
    authToken = null;
    localStorage.removeItem('token');
    showAuth();
});

function showAuth() {
    authScreen.classList.remove('hidden');
    dashboard.classList.add('hidden');
}

function showDashboard() {
    authScreen.classList.add('hidden');
    dashboard.classList.remove('hidden');
    loadServers();
}

// Navigation
menuItems.forEach(item => {
    item.addEventListener('click', () => {
        menuItems.forEach(i => i.classList.remove('active'));
        item.classList.add('active');
        const target = item.dataset.target;
        pages.forEach(p => {
            p.classList.remove('active');
            if (p.id === target) p.classList.add('active');
        });

        if (target === 'servers-page') loadServers();
        if (target === 'logs-page') {
            fetchLogs();
            startLogsAutoRefresh();
        } else {
            stopLogsAutoRefresh();
        }
        if (target === 'config-page') {
            loadGlobalConfig();
        }
    });
});

// Servers
async function loadServers() {
    try {
        const res = await fetchAuthenticated(`${API_BASE}/servers`);
        if (!res) return;
        servers = await res.json();
        renderServers();
    } catch (e) {
        console.error(e);
    }
}

function renderServers() {
    const list = document.getElementById('servers-list');
    list.innerHTML = '';
    servers.forEach(s => {
        const card = document.createElement('div');
        card.className = 'card server-card';
        card.innerHTML = `
            <h3>${s.Name}</h3>
            <div class="server-info"><i class="fa-solid fa-globe"></i> ${t('port')}: ${s.HTTPPort}</div>
            <div class="server-info"><i class="fa-solid fa-link"></i> ${s.EmbyHost}</div>
            <div class="server-info"><i class="fa-solid fa-folder"></i> ${s.MountPath}</div>
            <div class="server-actions">
                <button class="btn btn-sm btn-secondary" onclick="editServer(${s.ID})"><i class="fa-solid fa-pen"></i></button>
                <button class="btn btn-sm btn-danger" onclick="deleteServer(${s.ID})"><i class="fa-solid fa-trash"></i></button>
            </div>
        `;
        list.appendChild(card);
    });
}

// Server Modal
function showServerModal(id = null) {
    const modal = document.getElementById('server-modal');
    const form = document.getElementById('server-form');
    modal.classList.remove('hidden');

    if (id) {
        document.getElementById('modal-title').textContent = t('editServer');
        const s = servers.find(x => x.ID === id);
        if (s) {
            document.getElementById('server-id').value = s.ID;
            document.getElementById('server-name').value = s.Name;
            document.getElementById('server-port').value = s.HTTPPort;
            document.getElementById('server-emby-host').value = s.EmbyHost;
            document.getElementById('server-emby-token').value = s.EmbyToken || '';
            document.getElementById('server-mount-path').value = s.MountPath;
            document.getElementById('server-local-media-root').value = s.LocalMediaRoot || '';
            document.getElementById('server-openlist-host').value = s.OpenlistHost || '';
            document.getElementById('server-openlist-token').value = s.OpenlistToken || '';
            document.getElementById('server-internal-redirect').checked = s.InternalRedirectEnable;
            const dl = s.DirectLinkCacheExpired || '10m';
            document.getElementById('server-dl-cache-value').value = parseInt(dl) || 10;
            document.getElementById('server-dl-cache-unit').value = dl.replace(/\d+/,'') || 'm';
        }
    } else {
        document.getElementById('modal-title').textContent = t('addServer');
        form.reset();
        document.getElementById('server-id').value = '';
        // Explicitly set checkboxes to false
        document.getElementById('server-internal-redirect').checked = false;
    }
}

function closeServerModal() {
    document.getElementById('server-modal').classList.add('hidden');
}

window.editServer = showServerModal;
window.deleteServer = async (id) => {
    if (!confirm(t('deleteConfirm'))) return;
    await fetchAuthenticated(`${API_BASE}/servers/${id}`, { method: 'DELETE' });
    loadServers();
};

document.getElementById('server-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const id = document.getElementById('server-id').value;
    const data = {
        Name: document.getElementById('server-name').value,
        HTTPPort: parseInt(document.getElementById('server-port').value),
        EmbyHost: document.getElementById('server-emby-host').value,
        EmbyToken: document.getElementById('server-emby-token').value,
        MountPath: document.getElementById('server-mount-path').value,
        LocalMediaRoot: document.getElementById('server-local-media-root').value,
        OpenlistHost: document.getElementById('server-openlist-host').value,
        OpenlistToken: document.getElementById('server-openlist-token').value,
        InternalRedirectEnable: document.getElementById('server-internal-redirect').checked,
        DirectLinkCacheExpired: `${document.getElementById('server-dl-cache-value').value}${document.getElementById('server-dl-cache-unit').value}`,
        DisableProxy: false
    };

    const method = id ? 'PUT' : 'POST';
    const url = id ? `${API_BASE}/servers/${id}` : `${API_BASE}/servers`;

    const res = await fetchAuthenticated(url, {
        method: method,
        body: JSON.stringify(data),
        headers: { 'Content-Type': 'application/json' }
    });

    if (res) {
        closeServerModal();
        loadServers();
    }
});

// Logs
async function fetchLogs() {
    const filter = document.getElementById('log-filter').value;
    const serverFilter = document.getElementById('log-server-filter').value;
    
    // Update server list in filter if needed
    const serverSelect = document.getElementById('log-server-filter');
    // Keep the "all" option
    const allOption = serverSelect.querySelector('option[value="all"]');
    const currentVal = serverSelect.value;
    
    // Refresh options but keep selection
    serverSelect.innerHTML = '';
    serverSelect.appendChild(allOption);
    
    if (servers && servers.length > 0) {
        servers.forEach(s => {
            const opt = document.createElement('option');
            opt.value = s.Name;
            opt.textContent = s.Name;
            serverSelect.appendChild(opt);
        });
    }
    serverSelect.value = currentVal; // Restore selection

    const res = await fetchAuthenticated(`${API_BASE}/logs`);
    if (!res) return;
    let logs = await res.json();

    if (serverFilter !== 'all') {
        logs = logs.filter(l => {
            // Log format: [INFO] [ServerName] ...
            // We need to parse server name from message or adjust backend log format
            // Backend logs: [ServerName] Message...
            // Let's check how captureLogs works in manager.go:
            // logs.Info("[%s] %s", s.Name, strings.TrimSpace(line))
            // So the Message field starts with "[ServerName] "
            
            return l.Message.startsWith(`[${serverFilter}]`);
        });
    }

    if (filter === 'playback') {
        const reStream = /\/(emby\/)?videos\/.+\/(stream|universal|original)(\.\w+)?/i;
        const reM3u8 = /(master|main)\.m3u8/i;
        logs = logs.filter(l => {
            const msg = l.Message || "";
            // Ignore INFO level logs for playback filter
            if (l.Level === 'INFO') {
                // Allow INFO logs for m3u8 playlists if they contain specific keywords
                if (!reM3u8.test(msg) && !reStream.test(msg)) return false;
            }
            
            return (
                /重定向|Redirect/.test(msg) ||
                /直链|Direct|本地|Local|播放|Play/.test(msg) ||
                reStream.test(msg) ||
                reM3u8.test(msg)
            );
        });
    }
    if (filter === 'error') {
        logs = logs.filter(l => l.Level === 'ERROR');
    }

    if (logsCutoffMs !== null) {
        logs = logs.filter(l => {
            const ms = new Date(l.CreatedAt).getTime();
            return ms > logsCutoffMs;
        });
    }
    const times = logs.map(l => new Date(l.CreatedAt).getTime()).filter(n => !isNaN(n));
    if (times.length) {
        const latest = Math.max(...times);
        if (latest > logsLatestMs) logsLatestMs = latest;
    }

    const tbody = document.querySelector('#logs-table tbody');
    tbody.innerHTML = '';
    logs.forEach(l => {
        const tr = document.createElement('tr');
        
        // Extract server name for display column
        let serverName = "-";
        let message = l.Message;
        
        // Try to parse [ServerName] from message start
        const match = message.match(/^\[(.*?)\] (.*)/);
        if (match) {
            serverName = match[1];
            message = match[2];
        }

        tr.innerHTML = `
            <td>${new Date(l.CreatedAt).toLocaleString()}</td>
            <td>${serverName}</td>
            <td class="log-level-${l.Level.toLowerCase()}">${l.Level}</td>
            <td style="word-break: break-all;">${message}</td>
        `;
        tbody.appendChild(tr);
    });
}
function clearLogs() {
    const tbody = document.querySelector('#logs-table tbody');
    if (tbody) tbody.innerHTML = '';
    logsCutoffMs = logsLatestMs || Date.now();
}
let logsAutoTimer = null;
let logsAutoMs = 3000;
let logsCutoffMs = null;
let logsLatestMs = 0;
function startLogsAutoRefresh() {
    if (logsAutoTimer) clearInterval(logsAutoTimer);
    logsAutoTimer = setInterval(fetchLogs, logsAutoMs);
}
function stopLogsAutoRefresh() {
    if (logsAutoTimer) {
        clearInterval(logsAutoTimer);
        logsAutoTimer = null;
    }
}
const logFilterEl = document.getElementById('log-filter');
if (logFilterEl) {
    logFilterEl.addEventListener('change', fetchLogs);
}
const logServerFilterEl = document.getElementById('log-server-filter');
if (logServerFilterEl) {
    logServerFilterEl.addEventListener('change', fetchLogs);
}
const logAutoToggle = document.getElementById('log-autorefresh-toggle');
if (logAutoToggle) {
    logAutoToggle.addEventListener('click', () => {
        logsAutoMs = logsAutoMs === 3000 ? 1000 : 3000;
        const label = logsAutoMs === 1000 ? 'autoFast' : 'autoNormal';
        logAutoToggle.querySelector('span').textContent = t(label);
        if (document.querySelector('#logs-page').classList.contains('active')) {
            startLogsAutoRefresh();
        }
    });
}

// Global config form
async function loadGlobalConfig() {
    const res = await fetchAuthenticated(`${API_BASE}/global-config`);
    if (!res) return;
    const g = await res.json();
    document.getElementById('g-episodes').checked = !!g.EpisodesUnplayPrior;
    document.getElementById('g-resort').checked = !!g.ResortRandomItems;
    document.getElementById('g-proxy').value = g.ProxyErrorStrategy || 'origin';
    document.getElementById('g-images').value = g.ImagesQuality || 100;
    document.getElementById('g-download').value = g.DownloadStrategy || '403';
    document.getElementById('g-cache-enable').checked = !!g.CacheEnable;
    document.getElementById('g-cache-expired').value = g.CacheExpired || '1d';
    document.getElementById('g-cache-whitelist').value = g.CacheWhiteList || '';
    document.getElementById('g-vp-enable').checked = !!g.VideoPreviewEnable;
    document.getElementById('g-vp-containers').value = g.VideoPreviewContainers || 'mp4,mkv';
    document.getElementById('g-vp-ignore').value = g.VideoPreviewIgnoreTemplateIds || 'LD,SD';
    document.getElementById('g-path').value = (g.PathEmby2Openlist || '').replace(/,/g, '\n');
    document.getElementById('g-log-disable').checked = !!g.LogDisableColor;
    document.getElementById('config-page').dataset.gid = g.ID;
    document.getElementById('g-strm-path').value = (g.StrmPathMap || '');
    document.getElementById('g-ltg-enable').checked = !!g.LTGEnable;
    document.getElementById('g-ltg-ffmpeg').checked = !!g.LTGFFmpegEnable;
    document.getElementById('g-ltg-virtual').value = g.LTGVirtualContainers || 'mp4,mkv';
    document.getElementById('g-ltg-strm').value = g.LTGStrmContainers || 'ts';
    document.getElementById('g-ltg-music').value = g.LTGMusicContainers || 'mp3,flac';
    document.getElementById('g-ltg-auto-remove').value = g.LTGAutoRemoveMaxCount || 6000;
    document.getElementById('g-ltg-refresh').value = g.LTGRefreshInterval || 10;
    document.getElementById('g-ltg-scan').value = (g.LTGScanPrefixes || '');
    document.getElementById('g-ltg-ignore').value = g.LTGIgnoreContainers || 'jpg,jpeg,png,txt,nfo,md';
    document.getElementById('g-ltg-threads').value = g.LTGThreads || 8;
    document.getElementById('g-ssl-enable').checked = !!g.SslEnable;
    document.getElementById('g-ssl-single').checked = !!g.SslSinglePort;
    document.getElementById('g-ssl-key').value = g.SslKey || '';
    document.getElementById('g-ssl-crt').value = g.SslCrt || '';
}
async function saveGlobalConfig() {
    const payload = {
        ID: parseInt(document.getElementById('config-page').dataset.gid || '1'),
        EpisodesUnplayPrior: document.getElementById('g-episodes').checked,
        ResortRandomItems: document.getElementById('g-resort').checked,
        ProxyErrorStrategy: document.getElementById('g-proxy').value,
        ImagesQuality: parseInt(document.getElementById('g-images').value || '100'),
        DownloadStrategy: document.getElementById('g-download').value,
        CacheEnable: document.getElementById('g-cache-enable').checked,
        CacheExpired: document.getElementById('g-cache-expired').value,
        CacheWhiteList: document.getElementById('g-cache-whitelist').value.trim(),
        VideoPreviewEnable: document.getElementById('g-vp-enable').checked,
        VideoPreviewContainers: document.getElementById('g-vp-containers').value,
        VideoPreviewIgnoreTemplateIds: document.getElementById('g-vp-ignore').value,
        PathEmby2Openlist: document.getElementById('g-path').value.trim(),
        LogDisableColor: document.getElementById('g-log-disable').checked
    };
    payload.StrmPathMap = document.getElementById('g-strm-path').value.trim();
    payload.LTGEnable = document.getElementById('g-ltg-enable').checked;
    payload.LTGFFmpegEnable = document.getElementById('g-ltg-ffmpeg').checked;
    payload.LTGVirtualContainers = document.getElementById('g-ltg-virtual').value.trim();
    payload.LTGStrmContainers = document.getElementById('g-ltg-strm').value.trim();
    payload.LTGMusicContainers = document.getElementById('g-ltg-music').value.trim();
    payload.LTGAutoRemoveMaxCount = parseInt(document.getElementById('g-ltg-auto-remove').value || '6000');
    payload.LTGRefreshInterval = parseInt(document.getElementById('g-ltg-refresh').value || '10');
    payload.LTGScanPrefixes = document.getElementById('g-ltg-scan').value.trim();
    payload.LTGIgnoreContainers = document.getElementById('g-ltg-ignore').value.trim();
    payload.LTGThreads = parseInt(document.getElementById('g-ltg-threads').value || '8');
    payload.SslEnable = document.getElementById('g-ssl-enable').checked;
    payload.SslSinglePort = document.getElementById('g-ssl-single').checked;
    payload.SslKey = document.getElementById('g-ssl-key').value.trim();
    payload.SslCrt = document.getElementById('g-ssl-crt').value.trim();
    const res = await fetchAuthenticated(`${API_BASE}/global-config`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
    });
    if (res && res.ok) {
        alert(t('success'));
    }
}

// Helpers
async function fetchAuthenticated(url, options = {}) {
    if (!options.headers) options.headers = {};
    options.headers['Authorization'] = authToken;

    const res = await fetch(url, options);
    if (res.status === 401) {
        authToken = null;
        localStorage.removeItem('token');
        showAuth();
        return null;
    }
    return res;
}

// Start
init();
