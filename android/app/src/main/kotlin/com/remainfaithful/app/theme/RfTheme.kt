package com.remainfaithful.app.theme

import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.interaction.collectIsPressedAsState
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.RowScope
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.heightIn
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.Button
import androidx.compose.material3.ButtonDefaults
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedTextFieldDefaults
import androidx.compose.material3.Text
import androidx.compose.material3.darkColorScheme
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.graphicsLayer
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.font.FontFamily
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp

/**
 * Remain Faithful color system.
 * Change [gold], [goldLight], and [goldDeep] here (and the matching values in
 * res/values/colors.xml) when the gold moves to a brighter yellow.
 * Screens use these tokens. Do not hardcode brand colors elsewhere.
 */
object RfColors {
    const val navyArgb: Int = 0xFF12214C.toInt()
    const val goldArgb: Int = 0xFFD1AB4C.toInt()

    val navy = Color(navyArgb)
    val navyMid = Color(0xFF1C3066)
    val gold = Color(goldArgb)
    val goldLight = Color(0xFFF2D685)
    val goldDeep = Color(0xFFB88C2E)
    val danger = Color(0xFF9E2429)
    val success = Color(0xFF33C773)
    val info = Color(0xFF478FEF)
    val flag = Color(0xFFE6404D)
    val textPrimary = Color.White
    val textSecondary = Color.White.copy(alpha = 0.72f)
    val textTertiary = Color.White.copy(alpha = 0.60f)
    val card = Color.White.copy(alpha = 0.06f)
    val cardInset = Color.White.copy(alpha = 0.04f)
    val hairline = Color.White.copy(alpha = 0.08f)
    val stroke = Color.White.copy(alpha = 0.16f)
    val scrim = Color(0xC712214C)

    val page: Brush = Brush.verticalGradient(listOf(navy, navyMid))
    val goldFill: Brush = Brush.verticalGradient(listOf(goldLight, goldDeep))
}

object RfMetrics {
    val card = 16.dp
    val control = 14.dp
    val gutter = 20.dp
    val section = 24.dp
    val minTap = 48.dp
}

private val scheme = darkColorScheme(
    primary = RfColors.gold,
    onPrimary = RfColors.navy,
    secondary = RfColors.goldDeep,
    background = RfColors.navy,
    onBackground = RfColors.textPrimary,
    surface = RfColors.navyMid,
    onSurface = RfColors.textPrimary,
    error = RfColors.danger,
    onError = RfColors.textPrimary,
)

private val type = androidx.compose.material3.Typography(
    headlineLarge = TextStyle(fontFamily = FontFamily.Serif, fontWeight = FontWeight.Bold, fontSize = 34.sp, color = RfColors.textPrimary),
    headlineMedium = TextStyle(fontFamily = FontFamily.Serif, fontWeight = FontWeight.Bold, fontSize = 26.sp, color = RfColors.textPrimary),
    titleLarge = TextStyle(fontFamily = FontFamily.Serif, fontWeight = FontWeight.Bold, fontSize = 22.sp, color = RfColors.textPrimary),
    titleMedium = TextStyle(fontFamily = FontFamily.Serif, fontWeight = FontWeight.SemiBold, fontSize = 18.sp, color = RfColors.gold),
    bodyLarge = TextStyle(fontSize = 16.sp, color = RfColors.textPrimary),
    bodyMedium = TextStyle(fontSize = 15.sp, color = RfColors.textSecondary),
    labelSmall = TextStyle(fontSize = 11.sp, fontWeight = FontWeight.Bold, letterSpacing = 1.2.sp, color = RfColors.gold),
)

@Composable
fun RfTheme(content: @Composable () -> Unit) {
    MaterialTheme(colorScheme = scheme, typography = type, content = content)
}

fun Modifier.rfCard(inset: Boolean = false): Modifier {
    val shape = RoundedCornerShape(RfMetrics.card)
    return this
        .clip(shape)
        .background(if (inset) RfColors.cardInset else RfColors.card)
        .border(1.dp, RfColors.stroke, shape)
}

@Composable
fun rfFieldColors() = OutlinedTextFieldDefaults.colors(
    focusedBorderColor = RfColors.gold,
    unfocusedBorderColor = RfColors.hairline,
    focusedTextColor = RfColors.textPrimary,
    unfocusedTextColor = RfColors.textPrimary,
    cursorColor = RfColors.gold,
    focusedLabelColor = RfColors.gold,
    unfocusedLabelColor = RfColors.textTertiary,
    focusedContainerColor = RfColors.card,
    unfocusedContainerColor = RfColors.card,
)

@Composable
fun RfButton(
    title: String,
    onClick: () -> Unit,
    modifier: Modifier = Modifier,
    enabled: Boolean = true,
    destructive: Boolean = false,
    secondary: Boolean = false,
) {
    val interaction = remember { MutableInteractionSource() }
    val pressed by interaction.collectIsPressedAsState()
    val scale by animateFloatAsState(if (pressed && enabled) 0.96f else 1f, label = "press")
    val container = when {
        destructive -> RfColors.danger
        secondary -> Color.Transparent
        else -> RfColors.gold
    }
    val content = when {
        destructive -> RfColors.textPrimary
        secondary -> RfColors.gold
        else -> RfColors.navy
    }
    Button(
        onClick = onClick,
        enabled = enabled,
        interactionSource = interaction,
        modifier = modifier
            .fillMaxWidth()
            .heightIn(min = RfMetrics.minTap)
            .graphicsLayer {
                scaleX = scale
                scaleY = scale
            },
        shape = RoundedCornerShape(RfMetrics.control),
        contentPadding = PaddingValues(horizontal = 22.dp, vertical = 12.dp),
        colors = ButtonDefaults.buttonColors(
            containerColor = container,
            contentColor = content,
            disabledContainerColor = RfColors.gold.copy(alpha = 0.35f),
            disabledContentColor = RfColors.navy.copy(alpha = 0.6f),
        ),
        border = if (secondary) androidx.compose.foundation.BorderStroke(1.dp, RfColors.gold) else null,
    ) {
        Text(title, fontWeight = FontWeight.SemiBold)
    }
}

@Composable
fun RfTextButton(title: String, onClick: () -> Unit, modifier: Modifier = Modifier) {
    Text(
        title,
        color = RfColors.textPrimary,
        modifier = modifier
            .heightIn(min = RfMetrics.minTap)
            .clickable(onClick = onClick),
        fontSize = 16.sp,
    )
}

@Composable
fun SectionLabel(title: String) {
    Text(
        title,
        style = MaterialTheme.typography.labelSmall,
        modifier = Modifier.fillMaxWidth(),
    )
}

@Composable
fun RowScope.SpacerWeight() {
    androidx.compose.foundation.layout.Spacer(Modifier.weight(1f))
}
