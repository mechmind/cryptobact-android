package net.goandroid.cryptobact;

import android.app.Activity;
import android.content.Context;
import android.util.Log;
import android.os.Bundle;
import android.view.View;
import android.view.MotionEvent;
import android.opengl.GLSurfaceView;
import android.opengl.GLES20;
import android.net.wifi.WifiManager;
import android.net.wifi.WifiManager.MulticastLock;
import javax.microedition.khronos.egl.EGLConfig;
import javax.microedition.khronos.opengles.GL10;

import java.io.File;
import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.io.BufferedInputStream;
import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.io.DataInputStream;
import java.io.IOException;
import java.io.InputStream;

import java.lang.ProcessBuilder;
import java.lang.Process;
import java.lang.reflect.Method;

public class MainActivity extends Activity {
	private GLSurfaceView gl_view;
    private WifiManager wifi;
    private MulticastLock mLock;

    @Override public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        wifi = (WifiManager) getSystemService(Context.WIFI_SERVICE); 
        if(wifi != null){
            mLock = wifi.createMulticastLock("infektolock");
            mLock.acquire();
        }
		initGL();
	}

	private static void setupContextPreserve(GLSurfaceView gl_view) {
		try {
			Method setPreserveEGLContextOnPause_method = gl_view.getClass().getMethod("setPreserveEGLContextOnPause", new Class[]{Boolean.TYPE});
			setPreserveEGLContextOnPause_method.invoke(gl_view, new Object[]{Boolean.TRUE});
			Log.i("MainActivity", "Enabled context preservation");
		} catch (Exception e) {
			Log.i("MainActivity", "Failed to setup context preservation " + e.getMessage());
		}
	}

	private void initGL() {
		gl_view = new GLSurfaceView(this);
		setupContextPreserve(gl_view);
		gl_view.setEGLContextClientVersion(2);
		gl_view.setRenderer(new GLSurfaceView.Renderer() {
			@Override public void onDrawFrame(GL10 glUnused) {
				Engine.drawFrame();
			}
			@Override public void onSurfaceCreated(GL10 glUnused, EGLConfig config) {
				Engine.init();
			}
			@Override public void onSurfaceChanged(GL10 glUnused, int width, int height) {
				Engine.resize(width, height);
			}
		});
		gl_view.setOnTouchListener(new View.OnTouchListener() {
			@Override public boolean onTouch(View v, MotionEvent ev) {
				if (ev.getActionIndex() == 0) {
					Engine.onTouch(ev.getActionMasked(), ev.getX(), ev.getY());
					return true;
				}
				return false;
			}
		});
		setContentView(gl_view);
    }

	@Override protected void onResume() {
		super.onResume();
		gl_view.onResume();
        mLock.acquire();
	}

	@Override protected void onPause() {
		super.onPause();
		gl_view.onPause();
        mLock.release();
	}

	@Override public void onDestroy() {
		super.onDestroy();
	}
}
