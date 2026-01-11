package handlers

import (
	"net/http"
)

// ServeTestClient serves a simple HTML test client
func ServeTestClient(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Media Streaming Client</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        
        h1 {
            color: white;
            text-align: center;
            margin-bottom: 40px;
            font-size: 2.5rem;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
        }

        .source-selector {
            background: white;
            border-radius: 15px;
            padding: 20px;
            margin-bottom: 30px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.2);
            display: flex;
            gap: 15px;
            align-items: center;
            justify-content: center;
        }

        .source-selector label {
            font-weight: 600;
            color: #333;
        }

        .source-toggle {
            display: flex;
            gap: 10px;
        }

        .source-toggle button {
            padding: 10px 25px;
            border: 2px solid #667eea;
            background: white;
            color: #667eea;
            border-radius: 8px;
            cursor: pointer;
            font-weight: 600;
            transition: all 0.3s;
        }

        .source-toggle button.active {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
        }

        .source-toggle button:hover:not(.active) {
            background: #f0f0f0;
        }
        
        .section {
            background: white;
            border-radius: 15px;
            padding: 30px;
            margin-bottom: 30px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.2);
        }
        
        h2 {
            color: #667eea;
            margin-bottom: 20px;
            font-size: 1.8rem;
            display: flex;
            align-items: center;
            gap: 10px;
        }

        .source-badge {
            font-size: 0.7rem;
            padding: 4px 12px;
            border-radius: 12px;
            background: #e3f2fd;
            color: #1976d2;
            font-weight: 600;
        }
        
        .controls {
            display: flex;
            gap: 15px;
            margin-bottom: 20px;
            flex-wrap: wrap;
        }
        
        input {
            flex: 1;
            min-width: 200px;
            padding: 12px 20px;
            border: 2px solid #e0e0e0;
            border-radius: 8px;
            font-size: 16px;
            transition: border-color 0.3s;
        }
        
        input:focus {
            outline: none;
            border-color: #667eea;
        }
        
        button {
            padding: 12px 30px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border: none;
            border-radius: 8px;
            cursor: pointer;
            font-size: 16px;
            font-weight: 600;
            transition: transform 0.2s, box-shadow 0.2s;
        }
        
        button:hover {
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
        }
        
        button:active {
            transform: translateY(0);
        }
        
        audio {
            width: 100%;
            margin-top: 15px;
            border-radius: 8px;
        }
        
        .image-container {
            margin-top: 20px;
            text-align: center;
        }
        
        .image-container img {
            max-width: 100%;
            height: auto;
            border-radius: 10px;
            box-shadow: 0 5px 20px rgba(0,0,0,0.2);
        }
        
        .file-list {
            margin-top: 15px;
            padding: 15px;
            background: #f8f9fa;
            border-radius: 8px;
            max-height: 300px;
            overflow-y: auto;
        }
        
        .file-item {
            padding: 10px;
            margin: 5px 0;
            background: white;
            border-radius: 5px;
            cursor: pointer;
            transition: background 0.2s;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        
        .file-item:hover {
            background: #e3f2fd;
        }
        
        .file-name {
            font-weight: 500;
            color: #333;
        }
        
        .file-size {
            color: #666;
            font-size: 0.9rem;
        }
        
        .status {
            margin-top: 15px;
            padding: 12px;
            border-radius: 8px;
            font-weight: 500;
        }
        
        .status.success {
            background: #d4edda;
            color: #155724;
        }
        
        .status.error {
            background: #f8d7da;
            color: #721c24;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🎵 Media Streaming Client</h1>
        
      
        
        <div class="section">
            <h2>
                🎵 Music Streaming
                <span class="source-badge" id="musicSourceBadge">Local</span>
            </h2>
            <div class="controls">
                <input type="text" id="musicFile" placeholder="Enter music filename (e.g., song.mp3)">
                <button onclick="loadMusic()">Load Music</button>
                <button onclick="listMusic()">List Available</button>
            </div>
            <audio id="audioPlayer" controls></audio>
            <div id="musicList" class="file-list" style="display:none;"></div>
            <div id="musicStatus"></div>
        </div>
        
        <div class="section">
            <h2>
                🖼️ Image Streaming
                <span class="source-badge" id="imageSourceBadge">Local</span>
            </h2>
            <div class="controls">
                <input type="text" id="imageFile" placeholder="Enter image filename (e.g., photo.jpg)">
                <button onclick="loadImage()">Load Image</button>
                <button onclick="listImages()">List Available</button>
            </div>
            <div class="image-container">
                <img id="imageDisplay" style="display:none;" alt="Streamed image">
            </div>
            <div id="imageList" class="file-list" style="display:none;"></div>
            <div id="imageStatus"></div>
        </div>
    </div>
    
    <script>
    
        
        function switchSource(source) {
           
            
        
            document.getElementById('minioBtn').classList.toggle('active', source === 'minio');
            
            // Update badges
            const badgeText = source === 'MinIO';
            document.getElementById('musicSourceBadge').textContent = badgeText;
            document.getElementById('imageSourceBadge').textContent = badgeText;
            
            // Clear current displays
            document.getElementById('audioPlayer').src = '';
            document.getElementById('imageDisplay').style.display = 'none';
            document.getElementById('musicList').style.display = 'none';
            document.getElementById('imageList').style.display = 'none';
            document.getElementById('musicStatus').textContent = '';
            document.getElementById('imageStatus').textContent = '';
        }
        
        function getEndpointPrefix() {
            return '/api';
        }
        
        function loadMusic() {
            const filename = document.getElementById('musicFile').value;
            if (!filename) {
                //showStatus('musicStatus', 'Please enter a filename', 'error');
                return;
            }
            
            const prefix = getEndpointPrefix();
            const audioPlayer = document.getElementById('audioPlayer');
            audioPlayer.src = prefix + '/music/' + encodeURIComponent(filename);
            audioPlayer.load();
            //showStatus('musicStatus', 'Loading: ' + filename + ' from ' + 'MINIO', 'success');
        }
        
        function loadImage() {
            const filename = document.getElementById('imageFile').value;
            if (!filename) {
                showStatus('imageStatus', 'Please enter a filename', 'error');
                return;
            }
            
            const prefix = getEndpointPrefix();
            const imageDisplay = document.getElementById('imageDisplay');
            imageDisplay.src = prefix + '/images/' + encodeURIComponent(filename);
            imageDisplay.style.display = 'block';
            //showStatus('imageStatus', 'Loading: ' + filename + ' from ' + ’MINIO', 'success');
        }
        
        async function listMusic() {
            try {
                const prefix = getEndpointPrefix();
	
                const response = await fetch(prefix + '/music');
                const data = await response.json();
                displayFileList('musicList', data.files, 'music');
                //const source = data.source || 'minio';
                //showStatus('musicStatus', 'Found ' + data.files.length + ' music files in ' + source+ 'success');
            } catch (error) {
                //showStatus('musicStatus', 'Error listing files: ' + error.message, 'error');
            }
        }
        
        async function listImages() {
            try {
                const prefix = getEndpointPrefix();
                const response = await fetch(prefix + '/images');
                const data = await response.json();
                displayFileList('imageList', data.files, 'image');
                //const source = data.source || currentSource;
                //showStatus('imageStatus', 'Found ' + data.files.length + ' image files in ' + source, 'success');
            } catch (error) {
                //showStatus('imageStatus', 'Error listing files: ' + error.message, 'error');
            }
        }
        
        function displayFileList(elementId, files, type) {
            const listElement = document.getElementById(elementId);
            if (files.length === 0) {
                listElement.innerHTML = '<p style="text-align:center;color:#666;">No files found</p>';
                listElement.style.display = 'block';
                return;
            }
            
            listElement.innerHTML = files.map(file => {
                const sizeKB = (file.size / 1024).toFixed(2);
                return '<div class="file-item" onclick="selectFile(\'' + file.name + '\', \'' + type + '\')"><span class="file-name">' + file.name + '</span><span class="file-size">' + sizeKB + ' KB</span></div>';
            }).join('');
            listElement.style.display = 'block';
        }
        
        function selectFile(filename, type) {
            if (type === 'music') {
                document.getElementById('musicFile').value = filename;
                loadMusic();
            } else if (type === 'image') {
                document.getElementById('imageFile').value = filename;
                loadImage();
            }
        }
        
        function showStatus(elementId, message, type) {
            const statusElement = document.getElementById(elementId);
            statusElement.textContent = message;
            statusElement.className = 'status ' + type;
        }
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
