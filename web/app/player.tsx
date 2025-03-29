"use client";

import { useEffect, useRef, useState } from "react";
import Hls from "hls.js";

export default function Player() {
  const videoRef = useRef<HTMLVideoElement | null>(null);

  useEffect(() => {
    if (!videoRef.current) return;

    const hls = new Hls();

    hls.loadSource(
      "http://localhost:8080/video/It.Happened.One.Night.1934.2160p.4K.BluRay.x265.10bit.AAC5.1-[YTS.MX].mkv.m3u8",
    );
    hls.attachMedia(videoRef.current);

    return () => {
      hls.destroy();
    };
  }, []);

  return <video ref={videoRef} controls className="h-screen" />;
}
