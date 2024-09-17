// const prepareBLOB = (FILE, fileMetaData: fileMetadata) => {
//     const enc  = new TextEncoder(); // always utf-8, Uint8Array()
//     const buf1 = enc.encode('!');
//     const buf2 = enc.encode(JSON.stringify(fileMetaData));
//     const buf3 = enc.encode("\r\n\r\n");    
//     const buf4 = ";
    
//     let sendData = new Uint8Array(buf1.byteLength + buf2.byteLength + buf3.byteLength + buf4.byteLength);
//     sendData.set(new Uint8Array(buf1), 0);
//     sendData.set(new Uint8Array(buf2), buf1.byteLength);
//     sendData.set(new Uint8Array(buf3), buf1.byteLength + buf2.byteLength);
//     sendData.set(new Uint8Array(buf4), buf1.byteLength + buf2.byteLength + buf3.byteLength);
// }

