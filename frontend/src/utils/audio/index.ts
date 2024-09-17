let mediaRecorder: MediaRecorder
let audioBlobs: Blob[] = [];
let capturedStream: MediaStream

const startAudioRecording = async (recordedAudio: audio) => {
    try {
    const stream = await navigator.mediaDevices.getUserMedia({ audio: { echoCancellation: true }})
    var media = new MediaRecorder(stream, {mimeType: 'audio/webm'})
    let audioBlobs = [];
    capturedStream = stream;
    media.ondataavailable = async (evt) => {
          if(evt.data && evt.data.size > 0) {
            audioBlobs.push(evt.data);
            recordedAudio.blobs = audioBlobs
          }
      }
  
      media.start()
      return {media, stream: capturedStream}
    } catch (error) {
      console.log(error)
      return {media: undefined, stream: undefined}

    }
    
  }


  const stopAudioRecording = async (recordedMedia: audio) => {

    recordedMedia.recorder.onstop = (e) => {
        if(recordedMedia.blobs.length > 0) {
        const audio = document.createElement("audio");
        audio.controls = true;
        if (recordedMedia.stream.getTracks().length > 0) {
            recordedMedia.stream.getAudioTracks().forEach(track => track.stop());
        }


      };

    }
      return recordedMedia



    // let recordedAudioBlob: Blob = new Blob()
    //     await new Promise(resolve => {
    //         if (!mediaRecorder) {
    //           resolve(null);
    //           return;
    //         }
        
    //     mediaRecorder.addEventListener('stop', () => {
    //         const mimeType = mediaRecorder.mimeType;
    //         const audioBlob = new Blob(audioBlobs, { type: mimeType });
      
    //         if (capturedStream) {
    //           capturedStream.getTracks().forEach(track => track.stop());
    //         }
      
    //         resolve(audioBlob);
    //         recordedAudioBlob = audioBlob
    //         console.log(audioBlob)
    //       });
          
    //       mediaRecorder.stop();
          
    //     });
        
    //     return mediaRecorder
    
  }


  const playAudioRecording = (recording: Blob) => {
    if (recording) {
        const audio = new Audio();
        audio.src = URL.createObjectURL(recording);
        audio.play();
      }

  }
  

  
  
 

  export { startAudioRecording, stopAudioRecording, playAudioRecording}