package com.remainfaithful.app.ui

import android.provider.Settings
import androidx.compose.animation.core.Animatable
import androidx.compose.animation.core.CubicBezierEasing
import androidx.compose.animation.core.tween
import androidx.compose.foundation.Canvas
import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.BoxWithConstraints
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.WindowInsets
import androidx.compose.foundation.layout.asPaddingValues
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.statusBars
import androidx.compose.runtime.Composable
import androidx.compose.runtime.CompositionLocalProvider
import androidx.compose.runtime.DisposableEffect
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableFloatStateOf
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.runtime.staticCompositionLocalOf
import androidx.compose.runtime.withFrameNanos
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.semantics.clearAndSetSemantics
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.geometry.Size
import androidx.compose.ui.graphics.Path
import androidx.compose.ui.graphics.PathMeasure
import androidx.compose.ui.graphics.StrokeCap
import androidx.compose.ui.graphics.StrokeJoin
import androidx.compose.ui.graphics.drawscope.Stroke
import androidx.compose.ui.graphics.graphicsLayer
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.platform.LocalDensity
import androidx.compose.ui.text.font.FontFamily
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.lerp
import androidx.compose.ui.unit.sp
import androidx.compose.material3.Text
import androidx.lifecycle.Lifecycle
import androidx.lifecycle.LifecycleEventObserver
import androidx.lifecycle.compose.LocalLifecycleOwner
import com.remainfaithful.app.theme.RfColors
import com.remainfaithful.core.Copy
import kotlinx.coroutines.delay
import kotlinx.coroutines.isActive

data class LaunchState(
    val draw: Float,
    val glow: Float,
    val letters: Float,
    val travel: Float,
    val contentAlpha: Float,
    val reduceMotion: Boolean,
    val blocksInput: Boolean,
)

val LocalLaunch = staticCompositionLocalOf {
    LaunchState(
        draw = 1f,
        glow = 0f,
        letters = 1f,
        travel = 1f,
        contentAlpha = 1f,
        reduceMotion = false,
        blocksInput = false,
    )
}

@Composable
fun LaunchChrome(model: AppModel, content: @Composable () -> Unit) {
    val reduce = rememberReduceMotion()
    val resumed = rememberResumed()
    val launch = rememberLaunch(model, reduce, resumed)
    Box(Modifier.fillMaxSize().background(RfColors.page)) {
        if (!reduce && model.phase != Phase.App) {
            AmbientRings(paused = !resumed.value)
        }
        CompositionLocalProvider(LocalLaunch provides launch) {
            content()
            if (model.phase == Phase.Loading) {
                BoxWithConstraints(Modifier.fillMaxSize()) {
                    LaunchLockup(maxHeight = maxHeight, modifier = Modifier.align(Alignment.TopCenter))
                }
            }
            if (launch.blocksInput) {
                Box(
                    Modifier
                        .fillMaxSize()
                        .clearAndSetSemantics {}
                        .clickable(
                            interactionSource = remember { MutableInteractionSource() },
                            indication = null,
                        ) {},
                )
            }
        }
    }
}

@Composable
fun LaunchLockup(maxHeight: Dp, modifier: Modifier = Modifier) {
    val launch = LocalLaunch.current
    val travel = if (launch.reduceMotion) 1f else launch.travel
    val draw = if (launch.reduceMotion) 1f else launch.draw
    val mark = lerp(104.dp, 52.dp, travel)
    val topInset = WindowInsets.statusBars.asPaddingValues().calculateTopPadding()
    val top = lerp((maxHeight - mark - 44.dp) / 2, topInset + 8.dp, travel)
    val reveal = if (launch.reduceMotion) 1f else launch.letters
    val titleSize = (40f + (30f - 40f) * travel).sp
    Column(
        modifier
            .padding(top = top)
            .fillMaxWidth(),
        horizontalAlignment = Alignment.CenterHorizontally,
    ) {
        ShieldDraw(
            draw = draw,
            glow = launch.glow,
            modifier = Modifier.size(mark * 0.89f, mark),
        )
        Row {
            Copy.appName.forEachIndexed { index, ch ->
                val span = 1f / Copy.appName.length
                val start = index * span * 0.84f
                val shown = ((reveal - start) / span.coerceAtLeast(0.08f)).coerceIn(0f, 1f)
                Text(
                    if (ch == ' ') " " else ch.toString(),
                    color = RfColors.textPrimary.copy(alpha = shown),
                    fontFamily = FontFamily.Serif,
                    fontWeight = FontWeight.Bold,
                    fontSize = titleSize,
                )
            }
        }
    }
}

@Composable
fun Modifier.launchReveal(): Modifier {
    val launch = LocalLaunch.current
    val fade = remember { Animatable(if (launch.reduceMotion) 0f else 1f) }
    LaunchedEffect(launch.reduceMotion) {
        if (launch.reduceMotion) {
            fade.animateTo(1f, tween(450, easing = EaseInOut))
        }
    }
    val alpha = if (launch.reduceMotion) fade.value else launch.contentAlpha
    val density = LocalDensity.current
    return graphicsLayer {
        this.alpha = alpha
        if (!launch.reduceMotion) {
            translationY = with(density) { ((1f - alpha) * 18f).dp.toPx() }
        }
    }
}

private val EaseInOut = CubicBezierEasing(0.42f, 0f, 0.58f, 1f)

@Composable
private fun rememberReduceMotion(): Boolean {
    val context = LocalContext.current
    return remember {
        Settings.Global.getFloat(context.contentResolver, Settings.Global.ANIMATOR_DURATION_SCALE, 1f) == 0f
    }
}

@Composable
private fun rememberResumed(): androidx.compose.runtime.MutableState<Boolean> {
    val state = remember { mutableStateOf(true) }
    val owner = LocalLifecycleOwner.current
    DisposableEffect(owner) {
        val observer = LifecycleEventObserver { _, event ->
            when (event) {
                Lifecycle.Event.ON_RESUME -> state.value = true
                Lifecycle.Event.ON_PAUSE, Lifecycle.Event.ON_STOP -> state.value = false
                else -> Unit
            }
        }
        owner.lifecycle.addObserver(observer)
        onDispose { owner.lifecycle.removeObserver(observer) }
    }
    return state
}

@Composable
private fun rememberLaunch(
    model: AppModel,
    reduce: Boolean,
    resumed: androidx.compose.runtime.State<Boolean>,
): LaunchState {
    var draw by remember { mutableFloatStateOf(if (reduce) 1f else 0f) }
    var glow by remember { mutableFloatStateOf(0f) }
    var letters by remember { mutableFloatStateOf(if (reduce) 1f else 0f) }
    var travel by remember { mutableFloatStateOf(if (reduce) 1f else 0f) }
    var contentAlpha by remember { mutableFloatStateOf(if (reduce) 1f else 0f) }
    var skip by remember { mutableStateOf(false) }

    LaunchedEffect(model.phase) {
        if (model.phase == Phase.App) skip = true
    }

    LaunchedEffect(reduce, skip) {
        if (reduce || skip) {
            draw = 1f
            letters = 1f
            travel = 1f
            contentAlpha = 1f
            return@LaunchedEffect
        }
        val drawMs = 2400f
        val letterLead = 2150f
        val letterMs = 1400f
        val holdMs = 700f
        val travelMs = 1200f
        val total = letterLead + letterMs + holdMs + travelMs
        var elapsed = 0f
        var last = 0L
        while (isActive && elapsed < total && !skip) {
            val now = withFrameNanos { it }
            if (!resumed.value) {
                last = 0L
                delay(48)
                continue
            }
            if (last == 0L) {
                last = now
                continue
            }
            val dt = ((now - last) / 1_000_000f).coerceIn(0f, 50f)
            last = now
            elapsed += dt
            draw = ease(elapsed / drawMs)
            val bloom = if (elapsed > drawMs) ease((elapsed - drawMs) / 900f) else 0f
            val settle = ((elapsed - letterLead - 400f) / 1200f).coerceIn(0f, 1f)
            glow = bloom * (1f - settle * 0.65f)
            letters = if (elapsed > letterLead) ease((elapsed - letterLead) / letterMs) else 0f
            val moved = if (elapsed > letterLead + letterMs + holdMs) {
                ease((elapsed - letterLead - letterMs - holdMs) / travelMs)
            } else 0f
            travel = moved
            contentAlpha = moved
        }
        draw = 1f
        letters = 1f
        travel = 1f
        glow = 0.16f
        contentAlpha = 1f
    }

    val blocks = !reduce && !skip && model.phase != Phase.App && travel < 1f
    return LaunchState(draw, glow, letters, travel, contentAlpha, reduce, blocks)
}

private fun ease(t: Float): Float {
    val x = t.coerceIn(0f, 1f)
    return if (x < 0.5f) 2f * x * x else 1f - ((-2f * x + 2f) * (-2f * x + 2f)) / 2f
}

@Composable
private fun AmbientRings(paused: Boolean) {
    var time by remember { mutableFloatStateOf(0f) }
    LaunchedEffect(paused) {
        if (paused) return@LaunchedEffect
        var last = 0L
        while (isActive) {
            val now = withFrameNanos { it }
            if (last != 0L) {
                time += ((now - last) / 1_000_000f).coerceAtMost(90f)
            }
            last = now
            delay(80)
        }
    }
    Canvas(Modifier.fillMaxSize()) {
        val center = Offset(size.width / 2f, size.height * 0.42f)
        val maxR = maxOf(size.width, size.height) * 0.62f
        repeat(3) { index ->
            val shift = index / 3f
            val breathe = 0.5f + 0.5f * kotlin.math.sin(time / 1000f * 0.32f + shift * Math.PI.toFloat() * 2f)
            val radius = maxR * (0.34f + 0.22f * index) * (0.94f + 0.07f * breathe)
            drawCircle(
                color = RfColors.gold.copy(alpha = 0.035f + 0.045f * breathe),
                radius = radius,
                center = center,
                style = Stroke(width = 1.dp.toPx()),
            )
        }
    }
}

@Composable
private fun ShieldDraw(draw: Float, glow: Float, modifier: Modifier = Modifier) {
    Canvas(modifier) {
        val shield = shieldOutline(size)
        val check = shieldCheck(size)
        fun paint(path: Path, fraction: Float) {
            if (fraction <= 0f) return
            val measure = PathMeasure()
            measure.setPath(path, false)
            val segment = Path()
            measure.getSegment(0f, measure.length * fraction.coerceIn(0f, 1f), segment, true)
            if (glow > 0f) {
                drawPath(
                    segment,
                    color = RfColors.gold.copy(alpha = 0.28f * glow),
                    style = Stroke(width = 6.dp.toPx(), cap = StrokeCap.Round, join = StrokeJoin.Round),
                )
            }
            drawPath(
                segment,
                color = RfColors.gold,
                style = Stroke(width = 1.6.dp.toPx(), cap = StrokeCap.Round, join = StrokeJoin.Round),
            )
        }
        paint(shield, draw / 0.78f)
        paint(check, (draw - 0.72f) / 0.28f)
    }
}

private fun shieldOutline(size: Size): Path {
    val sx = size.width / 32f
    val sy = size.height / 36f
    fun x(v: Float) = v * sx
    fun y(v: Float) = v * sy
    return Path().apply {
        moveTo(x(16f), y(0f))
        lineTo(x(2f), y(6f))
        lineTo(x(2f), y(18f))
        cubicTo(x(2f), y(26.284f), x(8.268f), y(33.916f), x(16f), y(36f))
        cubicTo(x(23.732f), y(33.916f), x(30f), y(26.284f), x(30f), y(18f))
        lineTo(x(30f), y(6f))
        close()
    }
}

private fun shieldCheck(size: Size): Path {
    val sx = size.width / 32f
    val sy = size.height / 36f
    fun x(v: Float) = v * sx
    fun y(v: Float) = v * sy
    return Path().apply {
        moveTo(x(11f), y(18f))
        lineTo(x(14.5f), y(21.5f))
        lineTo(x(21f), y(14f))
    }
}
